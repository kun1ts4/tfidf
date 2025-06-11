package service

import (
	"math"
	"sort"
	"tfidf/internal/model"
)

func CalculateTFIDF(allDocs [][]string, docIndex int) (docStats []model.Word, collectionStats []model.Word) {
	docCount := len(allDocs)
	if docIndex >= docCount {
		return nil, nil
	}

	collectionWordFreq := make(map[string]int)
	collectionTotalWords := 0
	wordDocFrequency := make(map[string]int)

	docWordFreq := make(map[string]int)
	docTotalWords := len(allDocs[docIndex])

	// Собираем статистику по всем документам
	for _, doc := range allDocs {
		seenInDoc := make(map[string]bool)
		for _, word := range doc {
			collectionWordFreq[word]++
			collectionTotalWords++

			if !seenInDoc[word] {
				wordDocFrequency[word]++
				seenInDoc[word] = true
			}
		}
	}

	seenInTargetDoc := make(map[string]bool)
	for _, word := range allDocs[docIndex] {
		docWordFreq[word]++
		if !seenInTargetDoc[word] {
			seenInTargetDoc[word] = true
		}
	}

	for word, count := range docWordFreq {
		df := wordDocFrequency[word]
		idf := math.Log(float64(docCount) / float64(df))
		tf := float64(count) / float64(docTotalWords)

		docStats = append(docStats, model.Word{
			Word: word,
			TF:   tf,
			IDF:  idf,
			Freq: count,
		})
	}

	for word, count := range collectionWordFreq {
		df := wordDocFrequency[word]
		idf := math.Log(float64(docCount) / float64(df))
		tf := float64(count) / float64(collectionTotalWords)

		collectionStats = append(collectionStats, model.Word{
			Word: word,
			TF:   tf,
			IDF:  idf,
			Freq: count,
		})
	}

	return docStats, collectionStats
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
