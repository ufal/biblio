package biblio

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ufal/biblio/xmldb"
)

var (
	// ErrRecordNotAuthor is returned if the provided record doesn't represent an author.
	ErrRecordNotAuthor = errors.New("ill-typed record, expected 'author'")
	// ErrNoNameComponents is returned if the provided slice of name components is empty.
	ErrNoNameComponents = errors.New("expected at least one name component")
)

// Author is an author figuring in bibliographical entries.
type Author struct {
	// ID is the ID of the record.
	ID string `json:"id"`
	// Parent is the ID of the record's parent.
	Parent                string `json:"parent,omitempty"`
	FirstName             string `json:"first_name,omitempty"`
	MiddleName            string `json:"middle_name,omitempty"`
	LastName              string `json:"last_name,omitempty"`
	PrimaryEntry          string `json:"primary_entry,omitempty"`
	AcademicDegreesPrefix string `json:"academic_degrees_prefix,omitempty"`
	AcademicDegreesSuffix string `json:"academic_degrees_suffix,omitempty"`
	Nationality           string `json:"nationality,omitempty"`
	University            string `json:"university,omitempty"`
	Faculty               string `json:"faculty,omitempty"`
	Department            string `json:"department,omitempty"`
	CUniNumber            string `json:"cuni_number,omitempty"`
	CUniAffCode           string `json:"cuni_aff_code,omitempty"`
	Email                 string `json:"email,omitempty"`
	Note                  string `json:"note,omitempty"`
}

// NewAuthorFromNameComponents returns a new author with the provided name components.
func NewAuthorFromNameComponents(comps []string) (*Author, error) {
	author := &Author{ID: uuid.NewString()}
	switch len(comps) {
	case 0:
		return nil, ErrNoNameComponents
	case 1:
		author.FirstName = comps[0]
	case 2:
		author.FirstName = comps[0]
		author.LastName = comps[1]
	default:
		author.FirstName = comps[0]
		author.MiddleName = strings.Join(comps[1:len(comps)-1], " ")
		author.LastName = comps[len(comps)-1]
	}
	return author, nil
}

// NewAuthorFromRecord returns a new instance of the Author type based on the provided generic record.
func NewAuthorFromRecord(r *xmldb.Record) (*Author, error) {
	if r.Type != "author" {
		return nil, ErrRecordNotAuthor
	}
	a := &Author{ID: r.Id}
	if r.Parent != "0" {
		a.Parent = r.Parent
	}
	for _, field := range r.Fields {
		value := strings.TrimSpace(field.Value)
		switch field.Name {
		case "FirstName":
			a.FirstName = value
		case "MiddleName":
			a.MiddleName = value
		case "LastName":
			a.LastName = value
		case "PrimaryEntry":
			a.PrimaryEntry = value
		case "AcademicDegreesPrefix":
			a.AcademicDegreesPrefix = value
		case "AcademicDegreesSuffix":
			a.AcademicDegreesSuffix = value
		case "Nationality":
			a.Nationality = value
		case "University":
			a.University = value
		case "Faculty":
			a.Faculty = value
		case "Department":
			a.Department = value
		case "CUniNumber":
			a.CUniNumber = value
		case "CUniAffCode":
			a.CUniAffCode = value
		case "Email":
			a.Email = value
		case "Note":
			a.Note = value
		default:
			return nil, fmt.Errorf("unknown field '%s' (%s)", field.Name, field.Label)
		}
	}
	return a, nil
}

// FullName returns the full name of the author.
// It is the concatenation of the author's name components.
func (a *Author) FullName() string {
	comps := make([]string, 0, 3)
	if a.FirstName != "" {
		comps = append(comps, a.FirstName)
	}
	if a.MiddleName != "" {
		comps = append(comps, a.MiddleName)
	}
	if a.LastName != "" {
		comps = append(comps, a.LastName)
	}
	return strings.Join(comps, " ")
}
