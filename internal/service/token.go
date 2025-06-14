package service

import (
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

type TokenService struct {
	secret    string
	blacklist sync.Map
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
	if t.isBlacklisted(tokenString) {
		return "", jwt.ErrTokenInvalidClaims
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	return claims.Username, nil
}

func (t *TokenService) parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (t *TokenService) addToBlacklist(tokenString string, expiration time.Duration) {
	t.blacklist.Store(tokenString, true)
	time.AfterFunc(expiration, func() {
		t.blacklist.Delete(tokenString)
	})
}

func (t *TokenService) isBlacklisted(tokenString string) bool {
	_, ok := t.blacklist.Load(tokenString)
	return ok
}

func (t *TokenService) InvalidateToken(tokenString string) error {
	claims, err := t.parseToken(tokenString)
	if err != nil {
		return err
	}
	expiration := time.Until(claims.ExpiresAt.Time)
	t.addToBlacklist(tokenString, expiration)
	return nil
}
