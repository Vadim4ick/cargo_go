package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUsecase struct {
	secretAccess  []byte
	secretRefresh []byte
	accessExp     time.Duration
	refreshExp    time.Duration
}

func NewJWTService(secretAccess, secretRefresh string, accessExp, refreshExp time.Duration) *JwtUsecase {
	return &JwtUsecase{
		secretAccess:  []byte(secretAccess),
		secretRefresh: []byte(secretRefresh),
		accessExp:     accessExp,
		refreshExp:    refreshExp,
	}
}

// Генерация Access Token
func (j *JwtUsecase) GenerateAccess(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(j.accessExp).Unix(),
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretAccess)
}

// Генерация Refresh Token
func (j *JwtUsecase) GenerateRefresh(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(j.refreshExp).Unix(),
		"type": "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretRefresh)
}

// Валидация Access Token
func (j *JwtUsecase) ValidateAccess(tokenStr string) (string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secretAccess, nil
	})
	if err != nil || !t.Valid {
		return "", err
	}
	claims := t.Claims.(jwt.MapClaims)
	if claims["type"] != "access" {
		return "", errors.New("invalid token type")
	}
	return claims["sub"].(string), nil
}

// Валидация Refresh Token
func (j *JwtUsecase) ValidateRefresh(tokenStr string) (string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secretRefresh, nil
	})

	if err != nil || !t.Valid {
		return "", err
	}

	claims := t.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return "", errors.New("invalid token type")
	}
	return claims["sub"].(string), nil
}
