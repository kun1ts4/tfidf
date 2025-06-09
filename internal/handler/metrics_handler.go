package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tfidf/internal/model"
)

// GetMetrics godoc
// @Summary Get all metrics
// @Description Retrieve various metrics related to file processing
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} model.Metrics "Metrics retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /metrics [get]
func (h *Handler) GetMetrics(c *gin.Context) {
	peakUploadTime, err := h.repo.GetPeakUploadTime(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения пикового времени загрузки %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
		return
	}

	topFreqWords, err := h.repo.GetTopFreqWords(c.Request.Context(), 5)
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

	c.JSON(http.StatusOK, model.Metrics{
		PeakUploadTime:               peakUploadTime,
		TopFrequenciesWords:          topFreqWords,
		FilesProcessed:               filesProcessed,
		MinTimeProcessed:             minTimeProcessed,
		AvgTimeProcessed:             avgTimeProcessed,
		MaxTimeProcessed:             maxTimeProcessed,
		LatestFileProcessedTimestamp: latestFileProcessedTimestamp,
	})
}
