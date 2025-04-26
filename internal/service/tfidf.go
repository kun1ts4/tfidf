package service

import (
	"math"
	"tfidf/internal/model"
)

func CalculateTFIDF(allDocs [][]string) []model.WordTFIDF {
	docCount := len(allDocs)
	wordDocFrequency := make(map[string]int)
	wordTotalFrequency := make(map[string]int)
	totalWords := 0

	for _, doc := range allDocs {
		seen := make(map[string]bool)
		for _, word := range doc {
			wordTotalFrequency[word]++
			if !seen[word] {
				wordDocFrequency[word]++
				seen[word] = true
			}
		}
		totalWords += len(doc)
	}

	stats := make([]model.WordTFIDF, 0, len(wordTotalFrequency))

	for word, count := range wordTotalFrequency {
		tf := float64(count) / float64(totalWords)
		df := wordDocFrequency[word]
		idf := 0.0
		if df > 0 {
			idf = math.Log(float64(docCount) / float64(df))
		}

		stat := model.WordTFIDF{
			Word: word,
			TF:   tf,
			IDF:  idf,
		}
		stats = append(stats, stat)
	}
	return stats
}
