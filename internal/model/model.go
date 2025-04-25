package model

import "sort"

type WordTFIDF struct {
	Word string
	TF   float64
	IDF  float64
}

func TopIDFRange(all []WordTFIDF, n int, m int) []WordTFIDF {
	sort.Slice(all, func(i, j int) bool {
		return all[i].TF > all[j].TF
	})

	start := min(n, len(all))
	end := min(m, len(all))
	return all[start:end]
}
