package bibtex

import (
	"fmt"
)

// Entry is a BibTeX entry.
type Entry struct {
	// Type is the type of the entry.
	Type string `json:"type"`
	// ID is the ID of the entry.
	ID string `json:"id"`
	// AuthorIDs are the IDs of the authors.
	AuthorIDs []string `json:"author_ids,omitempty"`
	// Fields are the fields of the entry represented as a map.
	Fields map[string]string `json:"fields"`
}

// String returns a textual representation of the entry.
func (e *Entry) String() string {
	return fmt.Sprintf("%s %s %s", e.Type, e.ID, e.Fields)
}
