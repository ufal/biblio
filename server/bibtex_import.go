package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ufal/biblio"
	"github.com/ufal/biblio/bibdb"
	"github.com/ufal/biblio/bibtex"
	"github.com/ufal/biblio/bibtex/parser"
)

const (
	// MaxBibtexDataSize is the maximum size of a BibTeX string that can be imported.
	MaxBibtexDataSize = 1_000_000
)

type importBibtexResponse struct {
	Entries        []*bibtex.Entry  `json:"entries"`
	NumEntries     int              `json:"num_entries"`
	KnownAuthors   []*biblio.Author `json:"known_authors"`
	UnknownAuthors []*biblio.Author `json:"unknown_authors"`
}

type importBibtexHandler struct {
	store *bibdb.Store
}

func (h *importBibtexHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Printf("importBibtexHandler: unsupported method (%s)", req.Method)
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
		return
	}
	req.Body = http.MaxBytesReader(w, req.Body, MaxBibtexDataSize)
	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Print("importBibtexHandler: failed to read request body")
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	tokens := parser.Tokenise(string(b))
	rawEntries, err := parser.ParseEntries(tokens)
	if err != nil {
		log.Print("importBibtexHandler: failed to parse BibTeX data")
		http.Error(w, "failed to parse BibTeX data", http.StatusBadRequest)
		return
	}
	entries := make([]*bibtex.Entry, len(rawEntries))
	for i, rawEntry := range rawEntries {
		entry, err := rawEntry.Entry()
		if err != nil {
			log.Print("importBibtexHandler: failed to parse BibTeX entries")
			http.Error(w, "failed to parse BibTeX entries", http.StatusBadRequest)
			return
		}
		entries[i] = entry
	}
	authors, err := h.store.GetAllAuthors()
	if err != nil {
		log.Print("importBibtexHandler: failed to load all authors")
		http.Error(w, "failed to load all authors", http.StatusBadRequest)
		return
	}
	storedAuthors := make(map[string]*biblio.Author)
	for _, author := range authors {
		storedAuthors[author.FullName()] = author
	}
	importedAuthors := make(map[string]*biblio.Author)
	var knownAuthors, unknownAuthors []*biblio.Author
	for _, entry := range entries {
		if authors, ok := entry.Fields["author"]; ok {
			var authorIDs []string
			for _, fullName := range bibtex.ParseAuthors(authors) {
				if author, ok := importedAuthors[fullName]; ok {
					authorIDs = append(authorIDs, author.ID)
					continue
				}
				author, ok := storedAuthors[fullName]
				if ok {
					knownAuthors = append(knownAuthors, author)
					importedAuthors[fullName] = author
					authorIDs = append(authorIDs, author.ID)
				} else {
					comps := strings.Split(fullName, " ")
					author, err := biblio.NewAuthorFromNameComponents(comps)
					if err != nil {
						log.Print("importBibtexHandler: failed to create new author")
						http.Error(w, "failed to create new author", http.StatusBadRequest)
						return
					}
					unknownAuthors = append(unknownAuthors, author)
					importedAuthors[fullName] = author
					authorIDs = append(authorIDs, author.ID)
				}
			}
			entry.AuthorIDs = authorIDs
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	if err := enc.Encode(&importBibtexResponse{
		Entries:        entries,
		NumEntries:     len(entries),
		KnownAuthors:   knownAuthors,
		UnknownAuthors: unknownAuthors,
	}); err != nil {
		log.Print("importBibtexHandler: failed to write response")
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
