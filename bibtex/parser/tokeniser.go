package parser

import (
	"fmt"
	"unicode"
)

const (
	// Alphanum is an alphanumerical identifier.
	Alphanum = iota
	// Symbol is a one-character symbol.
	Symbol
	// EOF is the end-of-file token.
	EOF
)

// Token is a token serving as a basic building block entering the parser.
type Token struct {
	Text string
	Type int
	Line int
}

// String returns a textual representation of the token.
func (t *Token) String() string {
	return fmt.Sprintf("[%s]", t.Text)
}

// Tokenise split up a text into tokens.
func Tokenise(s string) []*Token {
	var (
		word   []rune
		tokens []*Token
		line   = 1
	)
	for _, r := range []rune(s) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			word = append(word, r)
		default:
			if len(word) > 0 {
				tokens = append(tokens, &Token{string(word), Alphanum, line})
				word = nil
			}
			tokens = append(tokens, &Token{string(r), Symbol, line})
			if r == '\n' {
				line++
			}
		}
	}
	if len(word) > 0 {
		tokens = append(tokens, &Token{string(word), Alphanum, line})
	}
	return append(tokens, &Token{"", EOF, line})
}
