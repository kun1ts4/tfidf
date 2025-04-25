package service

import (
	"tfidf/internal/model"
)

func CalculateTFIDF(words []string) []model.WordTFIDF {
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}

	uniqWords := len(counts)
	stats := make([]model.WordTFIDF, 0, uniqWords)

	for word, count := range counts {
		stat := model.WordTFIDF{
			Word: word,
			TF:   float64(count) / float64(uniqWords),
			IDF:  0,
		}
		stats = append(stats, stat)
	}

	return stats
}
