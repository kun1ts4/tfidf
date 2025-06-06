package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Auth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "токен должен быть заполнен"})
		c.Abort()
		return
	}

	tokenString := authHeader[len("Bearer "):]
	username, err := h.tokenService.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "неверный токен"})
		c.Abort()
		return
	}

	user, err := h.repo.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "неверный токен"})
		c.Abort()
		return
	}

	c.Set("userID", user.Id)

	ctx := context.WithValue(c.Request.Context(), "user", username)
	c.Request = c.Request.WithContext(ctx)
}
