package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tfidf/internal/config"
)

// Status godoc
// @Summary Status endpoint
// @Description app status check
// @Tags status
// @Success 200 {string} string "OK"
// @Router /status [get]
func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

// Version godoc
// @Summary Get app version
// @Description Returns current application version from config
// @Tags status
// @Produce json
// @Success 200 {object} map[string]interface{} "version"
// @Failure 500 {object} map[string]interface{} "error message"
// @Router /version [get]
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
