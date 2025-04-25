package main

import (
	"github.com/gin-gonic/gin"
	"tfidf/internal/handler"
)

func main() {
	r := gin.Default()

	r.Static("/styles", "./web/templates")
	r.LoadHTMLGlob("web/templates/*")

	r.GET("/ping", handler.Ping)
	r.GET("/", handler.ShowForm)
	r.POST("/upload", handler.UploadFile)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
