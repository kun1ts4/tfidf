package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tfidf/internal/config"
)

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
