package bibdb

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
	"github.com/ufal/biblio"
)

// Store is a persistent store for bibliographical records.
type Store struct {
	db *sql.DB
}

// NewStore opens a new store with the provided DSN.
func NewStore(dsn string) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	store := &Store{db}
	if err := store.setup(); err != nil {
		if err := store.Close(); err != nil {
			log.Print(err)
		}
		return nil, err
	}
	return store, nil
}

// Close closes the store.
func (s *Store) Close() error {
	return s.db.Close()
}

// GetAuthorByID returns the author with the given ID or nil if no such author exists.
func (s *Store) GetAuthorByID(id string) (*biblio.Author, error) {
	row := s.db.QueryRow("SELECT firstname, middlename, lastname FROM authors WHERE uid = $1", id)
	var firstName, middleName, lastName string
	if err := row.Scan(&firstName, &middleName, &lastName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &biblio.Author{
		ID:         id,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}, nil
}

// GetAllAuthors returns all authors.
func (s *Store) GetAllAuthors() ([]*biblio.Author, error) {
	rows, err := s.db.Query("SELECT uid, firstname, middlename, lastname FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var authors []*biblio.Author
	for rows.Next() {
		var id, firstName, middleName, lastName string
		if err := rows.Scan(&id, &firstName, &middleName, &lastName); err != nil {
			return nil, err
		}
		authors = append(authors, &biblio.Author{
			ID:         id,
			FirstName:  firstName,
			MiddleName: middleName,
			LastName:   lastName,
		})
	}
	if rows.NextResultSet() {
		log.Println("has more result sets")
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}

// InsertAuthor inserts a new author.
func (s *Store) InsertAuthor(author *biblio.Author) error {
	_, err := s.db.Exec("INSERT INTO authors (uid, firstname, middlename, lastname) VALUES ($1, $2, $3, $4)",
		author.ID, author.FirstName, author.MiddleName, author.LastName)
	return err
}
