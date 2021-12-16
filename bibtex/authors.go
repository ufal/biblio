package bibtex

import (
	"strings"
)

// ParseAuthors parses a list of authors into a slice of strings.
func ParseAuthors(s string) []string {
	var authors []string
	comps := strings.Split(s, "and")
	for _, el := range comps {
		el := strings.TrimSpace(el)
		if comps := strings.Split(el, ","); len(comps) > 1 {
			comps2 := make([]string, 0, len(comps))
			for _, el := range comps {
				comps2 = append(comps2, strings.TrimSpace(el))
			}
			el = strings.Join(comps2[1:], " ") + " " + comps2[0]
		}
		authors = append(authors, el)
	}
	return authors
}
