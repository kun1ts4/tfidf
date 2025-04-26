package parser

import (
	"regexp"
)

func ExtractDocuments(content []byte) []string {
	text := string(content)
	re := regexp.MustCompile(`(\r?\n){2,}`) // два или больше переводов строки подряд
	documents := re.Split(text, -1)
	return documents
}

func ExtractWords(text string) []string {
	re := regexp.MustCompile(`\p{L}+`)
	words := re.FindAllString(text, -1)
	return words
}
