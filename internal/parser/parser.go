package parser

import (
	"regexp"
)

func ExtractWords(content []byte) []string {
	text := string(content)

	re := regexp.MustCompile(`\p{L}+`)
	words := re.FindAllString(text, -1)
	return words
}
