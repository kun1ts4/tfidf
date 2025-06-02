package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) GetMetrics(c *gin.Context) {
	peakUploadTime, err := h.repo.GetPeakUploadTime(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения пикового времени загрузки %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	WordTFIDF, err := h.repo.GetTopFreqWords(c.Request.Context(), 5)
	if err != nil {
		log.Printf("ошибка получения топ-5 частотных слов %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	filesProcessed, err := h.repo.GetFilesProcessed(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения количества обработанных файлов %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	minTimeProcessed, err := h.repo.GetMinTimeProcessed(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения минимального времени обработки %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	avgTimeProcessed, err := h.repo.GetAvgTimeProcessed(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения среднего времени обработки %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	maxTimeProcessed, err := h.repo.GetMaxTimeProcessed(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения максимального времени обработки %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	latestFileProcessedTimestamp, err := h.repo.GetLatestFileProcessedTimestamp(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения времени последней обработки файла %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"peak_upload_time":                peakUploadTime,
		"top_frequencies_words":           WordTFIDF,
		"files_processed":                 filesProcessed,
		"min_time_processed":              minTimeProcessed,
		"avg_time_processed":              avgTimeProcessed,
		"max_time_processed":              maxTimeProcessed,
		"latest_file_processed_timestamp": latestFileProcessedTimestamp,
	})
}
