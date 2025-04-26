package model

import "sort"

type WordTFIDF struct {
	Word string
	TF   float64
	IDF  float64
}

func TopIDFRange(all []WordTFIDF, n int, m int) []WordTFIDF {
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
