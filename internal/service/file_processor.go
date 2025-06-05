package service

import (
	"tfidf/internal/model"
	"tfidf/internal/parser"
)

func ProcessFile(content []byte, top int) ([]model.Word, error) {
	documents := parser.ExtractDocuments(content)
	var allWords [][]string
	for _, document := range documents {
		allWords = append(allWords, parser.ExtractWords(document))
	}

	stats := CalculateTFIDF(allWords)
	res := TopIDFRange(stats, 0, top)
	return res, nil
}
