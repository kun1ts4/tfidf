package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"tfidf/internal/model"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account with username and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.User true "User credentials"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param user body model.User true "Login credentials"
// @Success 200 {object} model.TokenResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /login [post]
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

// Logout godoc
// @Summary User logout
// @Description Logout endpoint
// @Tags user
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} model.MessageResponse
// @Router /logout [get]
func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	tokenString := authHeader[len("Bearer "):]

	err := h.tokenService.InvalidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сделать токен недействительным"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "вы вышли из системы"})
}

// ChangeUserPassword godoc
// @Summary Change user password
// @Description Change password for authenticated user
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordRequest true "New password"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /user [patch]
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

// DeleteUser godoc
// @Summary Delete user account
// @Description Permanently delete authenticated user's account
// @Tags user
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} model.MessageResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /user [delete]
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
