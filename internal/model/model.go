package model

type Word struct {
	Word string  `json:"word"`
	TF   float64 `json:"-"`
	IDF  float64 `json:"-" `
	Freq int     `json:"frequency"`
}

type Metrics struct {
	PeakUploadTime               string  `json:"peak_upload_time"`
	TopFrequenciesWords          []Word  `json:"top_frequencies_words"`
	FilesProcessed               int     `json:"files_processed"`
	MinTimeProcessed             float64 `json:"min_time_processed"`
	AvgTimeProcessed             float64 `json:"avg_time_processed"`
	MaxTimeProcessed             float64 `json:"max_time_processed"`
	LatestFileProcessedTimestamp string  `json:"latest_file_processed_timestamp"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Document struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	AuthorId      int      `json:"author"`
	Collections   []string `json:"collections"`
	TimeProcessed float64  `json:"time_processed"`
}

type Stats struct {
	Words []Word `json:"words"`
}
