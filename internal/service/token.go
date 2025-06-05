package service

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenService struct {
	secret string
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{secret: secret}
}

func (t *TokenService) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(100 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", nil
	}
	return tokenString, err
}

func (t *TokenService) ValidateToken(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Username, nil
}
