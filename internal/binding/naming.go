package binding

import (
	"strings"
	"unicode"
)

func toSnake(s string) string {
	return strings.Join(parseWords(s), "_")
}

func parseWords(s string) (words []string) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}
	rs := []rune(s)
	word := ""
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if r == '_' || r == '-' || r == ' ' {
			if word != "" {
				words = append(words, word)
			}
			word = ""
			continue
		}
		if unicode.IsUpper(r) && ((i-1 > 0 && unicode.IsLower(rs[i-1])) || (i+1 < len(rs) && unicode.IsLower(rs[i+1]))) {
			if word != "" {
				words = append(words, word)
			}
			word = string(r)
			continue
		}
		word += string(r)
	}
	if word != "" {
		words = append(words, word)
	}
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return words
}
