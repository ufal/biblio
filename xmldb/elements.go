package xmldb

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// File is an XML file with a list of records.
type File struct {
	XMLName xml.Name  `xml:"Records"`
	Records []*Record `xml:"Record"`
}

// Record is a record with a list of fields.
type Record struct {
	XMLName xml.Name `xml:"Record"`
	Id      string   `xml:"Id,attr"`
	Parent  string   `xml:"Parent,attr"`
	Type    string   `xml:"Type,attr"`
	Fields  []*Field `xml:"Field"`
}

// String returns a textual representation of the record.
func (r *Record) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s (%s/%s)", r.Type, r.Id, r.Parent))
	for _, f := range r.Fields {
		sb.WriteString(fmt.Sprintf(" %s='%s'", f.Name, f.Value))
	}
	return sb.String()
}

// Field is a named field with a value and a label.
type Field struct {
	XMLName xml.Name `xml:"Field"`
	Name    string   `xml:"Name,attr"`
	Label   string   `xml:"Label,attr"`
	Value   string   `xml:",chardata"`
}

// LoadFile loads an XML file with records.
func LoadFile(r io.Reader) (*File, error) {
	f := new(File)
	if err := xml.NewDecoder(r).Decode(f); err != nil {
		return nil, err
	}
	return f, nil
}
