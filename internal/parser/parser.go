package parser

import (
	"regexp"
)

func ExtractWords(text string) []string {
	re := regexp.MustCompile(`\p{L}+`)
	words := re.FindAllString(text, -1)
	return words
}
