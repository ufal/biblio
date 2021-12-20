package tex

import (
	"fmt"
)

type textDirective struct {
	cmd string
	arg string
}

var (
	textDirectives = map[textDirective]string{
		// caron
		{"v", "c"}: "č",
		{"v", "C"}: "Č",
		{"v", "d"}: "ď",
		{"v", "D"}: "Ď",
		{"v", "e"}: "ě",
		{"v", "E"}: "Ě",
		{"v", "l"}: "ľ",
		{"v", "L"}: "Ľ",
		{"v", "n"}: "ň",
		{"v", "N"}: "Ň",
		{"v", "r"}: "ř",
		{"v", "R"}: "Ř",
		{"v", "s"}: "š",
		{"v", "S"}: "Š",
		{"v", "t"}: "ť",
		{"v", "T"}: "Ť",
		{"v", "z"}: "ž",
		{"v", "Z"}: "Ž",
		// ring
		{"r", "a"}: "å",
		{"r", "A"}: "Å",
		{"r", "u"}: "ů",
		{"r", "U"}: "Ů",
		// cedille
		{"c", "c"}: "ç",
		{"c", "C"}: "Ç",
		// umlaut
		{"\"", "a"}:   "ä",
		{"\"", "A"}:   "Ä",
		{"\"", "e"}:   "ë",
		{"\"", "E"}:   "Ë",
		{"\"", "\\i"}: "ï",
		{"\"", "I"}:   "Ï",
		{"\"", "o"}:   "ö",
		{"\"", "O"}:   "Ö",
		{"\"", "u"}:   "ü",
		{"\"", "U"}:   "Ü",
		{"\"", "y"}:   "ÿ",
		{"\"", "Y"}:   "Ÿ",
		// acute
		{"'", "a"}:   "á",
		{"'", "A"}:   "Á",
		{"'", "c"}:   "ć",
		{"'", "C"}:   "Ć",
		{"'", "e"}:   "é",
		{"'", "E"}:   "É",
		{"'", "\\i"}: "í",
		{"'", "I"}:   "Í",
		{"'", "l"}:   "ĺ",
		{"'", "L"}:   "Ĺ",
		{"'", "n"}:   "ń",
		{"'", "N"}:   "Ń",
		{"'", "o"}:   "ó",
		{"'", "O"}:   "Ó",
		{"'", "r"}:   "ŕ",
		{"'", "R"}:   "Ŕ",
		{"'", "s"}:   "ś",
		{"'", "S"}:   "Ś",
		{"'", "u"}:   "ú",
		{"'", "U"}:   "Ú",
		{"'", "y"}:   "ý",
		{"'", "Y"}:   "Ý",
		{"'", "z"}:   "ź",
		{"'", "Z"}:   "Ź",
		// grave
		{"`", "a"}:   "à",
		{"`", "A"}:   "À",
		{"`", "e"}:   "è",
		{"`", "E"}:   "È",
		{"`", "\\i"}: "ì",
		{"`", "I"}:   "Ì",
		{"`", "o"}:   "ò",
		{"`", "O"}:   "Ò",
		{"`", "u"}:   "ù",
		{"`", "U"}:   "Ù",
	}
)

// StringFromTex converts a TeX command representing a diacritic into a string.
func StringFromTex(cmd, arg string) (string, error) {
	if r, ok := textDirectives[textDirective{cmd, arg}]; ok {
		return r, nil
	} else {
		return "", fmt.Errorf("unknown command \\%s{%s}", cmd, arg)
	}
}
