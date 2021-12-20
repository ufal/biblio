package parser

import (
	"fmt"
	"strings"

	"github.com/ufal/biblio/bibtex"
	"github.com/ufal/biblio/encoding/tex"
)

// Commands is a (La)TeX command.
type Command struct {
	Name string
}

// String returns a textual representation of the command.
func (c *Command) String() string { return "\\" + c.Name }

// Block is a (La)TeX block.
type Block struct {
	Opening  string
	Closing  string
	Contents []interface{}
}

// String returns a textual representation of the block.
func (b *Block) String() string {
	var sb strings.Builder
	sb.WriteString(b.Opening)
	for _, el := range b.Contents {
		switch el := el.(type) {
		case string:
			sb.WriteString(el)
		case fmt.Stringer:
			sb.WriteString(el.String())
		default:
			sb.WriteString(fmt.Sprintf("(?%v?)", el))
		}
	}
	sb.WriteString(b.Closing)
	return sb.String()
}

// PlainText returns a plain text representation of the block.
func (b *Block) PlainText() (string, error) {
	var sb strings.Builder
	for i := 0; i < len(b.Contents); i++ {
		switch el := b.Contents[i].(type) {
		case string:
			sb.WriteString(el)
		case *Command:
			j := i + 1
			if j < len(b.Contents) {
				arg := b.Contents[j]
				if b, ok := arg.(*Block); ok {
					var err error
					arg, err = b.PlainText()
					if err != nil {
						return "", err
					}
				}
				if arg, ok := arg.(string); ok {
					if strings.TrimSpace(arg) == "" {
						sb.WriteString("\\" + el.Name)
						sb.WriteString(arg)
					} else {
						s, err := tex.StringFromTex(el.Name, arg)
						if err != nil {
							return "", err
						}
						sb.WriteString(s)
					}
					i++
					continue
				}
			}
			sb.WriteString(el.String())
			continue
		case *Block:
			pt, err := el.PlainText()
			if err != nil {
				return "", err
			}
			sb.WriteString(pt)
		case fmt.Stringer:
			sb.WriteString(el.String())
		default:
			sb.WriteString(fmt.Sprintf("(?%v?)", el))
		}
	}
	return sb.String(), nil
}

// RawEntry is a BibTeX entry with raw (unresolved) contents.
type RawEntry struct {
	Type   string
	Id     string
	Fields map[string]interface{}
}

// Entry returns a BibTeX entry.
func (e *RawEntry) Entry() (*bibtex.Entry, error) {
	fields := make(map[string]string)
	for key, val := range e.Fields {
		var strval string
		switch val := val.(type) {
		case string:
			strval = val
		case *Block:
			var err error
			strval, err = val.PlainText()
			if err != nil {
				return nil, err
			}
			fields[key] = strval
		}
	}
	return &bibtex.Entry{Type: e.Type, ID: e.Id, Fields: fields}, nil
}

func parseEntryType(tokens []*Token) (string, []*Token, error) {
	token := tokens[0]
	if token.Text != "@" {
		return "", nil, fmt.Errorf("expected '@' on line %d", token.Line)
	}
	tokens = tokens[1:]
	token = tokens[0]
	if token.Type != Alphanum {
		return "", nil, fmt.Errorf("expected alphanumeric identifier on line %d", token.Line)
	}
	return strings.ToLower(token.Text), tokens[1:], nil
}

func parseBlock(tokens []*Token) (*Block, []*Token, error) {
	token, sentinel := tokens[0], ""
	switch token.Text {
	case "{":
		sentinel = "}"
	case "\"":
		sentinel = "\""
	default:
		return nil, nil, fmt.Errorf("expected '{', '\"' on line %d", token.Line)
	}
	opening := token.Text
	var contents []interface{}
	tokens = tokens[1:]
	token = tokens[0]
	for {
		switch token.Text {
		case sentinel:
			return &Block{opening, sentinel, contents}, tokens[1:], nil
		case "{", "\"":
			var (
				block *Block
				err   error
			)
			block, tokens, err = parseBlock(tokens)
			if err != nil {
				return nil, nil, err
			}
			contents = append(contents, block)
			token = tokens[0]
			continue
		case "\\":
			tokens = tokens[1:]
			token = tokens[0]
			if token.Text == "" {
				return nil, nil, fmt.Errorf("unexpected EOF on line %d", token.Line)
			}
			contents = append(contents, &Command{token.Text})
		case "":
			return nil, nil, fmt.Errorf("unexpected EOF on line %d", token.Line)
		default:
			contents = append(contents, token.Text)
		}
		tokens = tokens[1:]
		token = tokens[0]
	}
}

func parseEntryId(tokens []*Token) (string, []*Token, error) {
	token := tokens[0]
	if token.Text == "," {
		return "", nil, fmt.Errorf("unexpected ',' on line %d", token.Line)
	}
	var id string
	for {
		if token.Text == "" {
			return "", nil, fmt.Errorf("unexpected EOF on line %d", token.Line)
		}
		id += token.Text
		tokens = tokens[1:]
		token = tokens[0]
		if token.Text == "," {
			return id, tokens[1:], nil
		}
	}
}

func parseEntryField(tokens []*Token) (string, interface{}, bool, []*Token, error) {
	tokens = eatWhiteSpacesAndComments(tokens)
	token := tokens[0]
	if token.Type != Alphanum {
		return "", "", false, nil, fmt.Errorf("expected alphanumeric identifier on line %d", token.Line)
	}
	key := strings.ToLower(token.Text)
	tokens = tokens[1:]
	tokens = eatWhiteSpacesAndComments(tokens)
	token = tokens[0]
	if token.Text != "=" {
		return "", "", false, nil, fmt.Errorf("expected '=' on line %d", tokens[0].Line)
	}
	tokens = tokens[1:]
	tokens = eatWhiteSpacesAndComments(tokens)
	token = tokens[0]
	if token.Text == "" {
		return "", "", false, nil, fmt.Errorf("unexpected EOF on line %d", tokens[0].Line)
	}
	var value interface{}
	if token.Type == Alphanum {
		value = token.Text
		tokens = tokens[1:]
	} else {
		var err error
		value, tokens, err = parseBlock(tokens)
		if err != nil {
			return "", "", false, nil, err
		}
	}
	tokens = eatWhiteSpacesAndComments(tokens)
	token = tokens[0]
	if token.Text == "" {
		return "", "", false, nil, fmt.Errorf("unexpected EOF on line %d", tokens[0].Line)
	}
	if token.Text != "," && token.Text != "}" {
		return "", "", false, nil, fmt.Errorf("expected ',', '}' on line %d", tokens[0].Line)
	}
	return key, value, token.Text == ",", tokens[1:], nil
}

func parseEntry(tokens []*Token) (*RawEntry, []*Token, error) {
	tokens = eatWhiteSpacesAndComments(tokens)
	typ, tokens, err := parseEntryType(tokens)
	if err != nil {
		return nil, nil, err
	}
	tokens = eatWhiteSpacesAndComments(tokens)
	if tokens[0].Text != "{" {
		return nil, nil, fmt.Errorf("expected '{' on line %d", tokens[0].Line)
	}
	id, tokens, err := parseEntryId(tokens[1:])
	if err != nil {
		return nil, nil, err
	}
	id = strings.TrimSpace(id)
	m := make(map[string]interface{})
	cont := true
	for cont {
		var (
			key   string
			value interface{}
			err   error
		)
		tokens = eatWhiteSpacesAndComments(tokens)
		if tokens[0].Text == "}" {
			tokens = tokens[1:]
			break
		}
		key, value, cont, tokens, err = parseEntryField(tokens)
		if err != nil {
			return nil, nil, err
		}
		m[key] = value
	}
	return &RawEntry{typ, id, m}, tokens, nil
}

// ParseEntries parses a sequence of tokens into a slice of raw BibTeX entries.
func ParseEntries(tokens []*Token) ([]*RawEntry, error) {
	var entries []*RawEntry
	for {
		tokens = eatWhiteSpacesAndComments(tokens)
		token := tokens[0]
		if token.Text == "" {
			return entries, nil
		}
		var (
			entry *RawEntry
			err   error
		)
		entry, tokens, err = parseEntry(tokens)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
}

func eatWhiteSpacesAndComments(tokens []*Token) []*Token {
	token := tokens[0]
	if token.Text == "%" {
	loop:
		for {
			tokens = tokens[1:]
			switch tokens[0].Text {
			case "":
				return tokens
			case "\n":
				break loop
			}
		}
		token = tokens[0]
	}
	if token.Text != "" && strings.TrimSpace(token.Text) == "" {
		return eatWhiteSpacesAndComments(tokens[1:])
	}
	return tokens
}
