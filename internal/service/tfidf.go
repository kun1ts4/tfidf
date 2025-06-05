package service

import (
	"math"
	"sort"
	"tfidf/internal/model"
)

func CalculateTFIDF(allDocs [][]string) []model.Word {
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

	stats := make([]model.Word, 0, len(wordTotalFrequency))

	for word, count := range wordTotalFrequency {
		tf := float64(count) / float64(totalWords)
		df := wordDocFrequency[word]
		idf := 0.0
		if df > 0 {
			idf = math.Log(float64(docCount) / float64(df))
		}

		stat := model.Word{
			Word: word,
			TF:   tf,
			IDF:  idf,
			Freq: count,
		}
		stats = append(stats, stat)
	}
	return stats
}

func TopIDFRange(all []model.Word, n int, m int) []model.Word {
	sort.Slice(all, func(i, j int) bool {
		if all[i].IDF == all[j].IDF {
			if all[i].TF == all[j].TF {
				return all[i].Word < all[j].Word
			}
			return all[i].TF > all[j].TF
		}
		return all[i].IDF > all[j].IDF
	})

	start := min(n, len(all))
	end := min(m, len(all))
	return all[start:end]
}
