package handler

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"tfidf/internal/config"
	"tfidf/internal/repository"
	"tfidf/internal/service"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func ShowForm(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ошибка загрузки файла",
		})
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ошибка чтения файла",
		})
	}

	top50, err := service.ProcessFile(content, 50)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	ctx := c.Request.Context()
	for _, word := range top50 {
		err = h.repo.RecordWordFrequency(ctx, word.Word, word.Freq)
		if err != nil {
			log.Printf("ошибка сохранения данных в базе %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "внутренняя ошибка сервера",
			})
			return
		}
	}

	fileName := fileHeader.Filename
	err = h.repo.RecordFileUpload(ctx, fileName)
	if err != nil {
		log.Printf("ошибка записи информации о загрузке файла %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "внутренняя ошибка сервера",
		})
		return
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"Words": top50,
	})
}

func (h *Handler) GetMetrics(c *gin.Context) {
	peakUploadTime, err := h.repo.GetPeakUploadTime(c.Request.Context())
	if err != nil {
		log.Printf("ошибка получения пикового времени загрузки %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка получения метрик",
		})
	}

	WordTFIDF, err := h.repo.GetTopFreqWords(c.Request.Context(), 5)

	c.JSON(http.StatusOK, gin.H{
		"peak_upload_time":      peakUploadTime,
		"top_frequencies_words": WordTFIDF,
	})
}

func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func Version(c *gin.Context) {
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "ошибка загрузки конфигурации",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"version": cfg.App.Version,
	})
}
