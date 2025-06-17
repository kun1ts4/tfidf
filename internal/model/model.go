package model

// Word represents a word with its frequency and metrics.
// @Description Represents a word and its associated metrics.
type Word struct {
	Word string  `json:"word"`
	TF   float64 `json:"tf"`
	IDF  float64 `json:"idf" `
	Freq int     `json:"-"`
}

// Metrics represents application metrics for processing files.
// @Description Contains various metrics related to file processing and word frequency analysis.
type Metrics struct {
	PeakUploadTime               string  `json:"peak_upload_time"`
	UsersCount                   int     `json:"registered_users_count"`
	FilesProcessed               int     `json:"files_processed"`
	MinTimeProcessed             float64 `json:"min_time_processed"`
	AvgTimeProcessed             float64 `json:"avg_time_processed"`
	MaxTimeProcessed             float64 `json:"max_time_processed"`
	LatestFileProcessedTimestamp string  `json:"latest_file_processed_timestamp"`
}

// User represents a system user.
// @Description Contains user credentials and identification information.
type User struct {
	Id       int    `json:"id" swaggerignore:"true"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Document represents a processed document.
// @Description Contains metadata and processing details of a document.
type Document struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	AuthorId      int      `json:"author"`
	Collections   []string `json:"collections"`
	TimeProcessed float64  `json:"time_processed"`
}

// Stats represents statistical data.
// @Description Contains aggregated word statistics.
type Stats struct {
	Words []Word `json:"words"`
}

// MessageResponse represents a response containing a message.
// @Description Contains a message for the user.
// @Example {"message": "пользователь успешно удален"}
type MessageResponse struct {
	Message string `json:"message" example:"пользователь успешно удален"`
}

// TokenResponse represents a response containing a JWT token.
// @Description Contains a JWT token for authentication.
// @Example {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ErrorResponse represents a response containing an error message.
// @Description Contains an error message for the user.
// @Example {"message": "неверный токен"}
type ErrorResponse struct {
	Message string `json:"message" example:"неверный токен"`
}

// ChangePasswordRequest represents a request to change a user's password.
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password"`
}
