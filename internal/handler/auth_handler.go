package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tfidf/internal/model"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный ввод"})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "требуются имя пользователя и пароль."})
		return
	}

	err := h.repo.CreateUser(c.Request.Context(), user)
	if err != nil {
		log.Printf("ошибка регистрации пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "не удалось зарегистрировать пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "пользователь успешно зарегистрирован"})
}

func (h *Handler) Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный ввод"})
		return
	}

	isPasswordCorrect, err := h.repo.CheckUserPassword(context.Background(), user)
	if err != nil {
		log.Printf("ошибка проверки пароля пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "внутренняя ошибка сервера"})
		return
	}

	if !isPasswordCorrect {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "неверный пароль"})
		return
	}

	tokenString, err := h.tokenService.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *Handler) Logout(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "вы вышли из системы"})
}

func (h *Handler) ChangeUserPassword(c *gin.Context) {
	var request struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "неверный ввод"})
		return
	}

	username, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "неверный токен"})
		return
	}

	err := h.repo.ChangeUserPassword(c.Request.Context(), username.(string), request.NewPassword)
	if err != nil {
		log.Printf("ошибка изменения пароля пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "пароль успешно изменен"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	username, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "неверный токен"})
		return
	}

	err := h.repo.DeleteUser(c.Request.Context(), username.(string))
	if err != nil {
		log.Printf("ошибка удаления пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "внутренняя ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "пользователь успешно удален"})
}
