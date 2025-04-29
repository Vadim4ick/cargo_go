package usecase

import (
	"errors"
	"test-project/internal/domain/user"
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
func (j *JwtUsecase) GenerateAccess(userID string, role user.Role) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(j.accessExp).Unix(),
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretAccess)
}

// Генерация Refresh Token
func (j *JwtUsecase) GenerateRefresh(userID string, role user.Role) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(j.refreshExp).Unix(),
		"type": "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretRefresh)
}

// Генерация Invite Token
func (j *JwtUsecase) GenerateInvite(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(5 * time.Minute).Unix(),
		"type":  "invite",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretAccess) // Можно использовать тот же secretAccess
}

// Валидация Access Token
func (j *JwtUsecase) ValidateAccess(tokenStr string) (string, user.Role, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secretAccess, nil
	})
	if err != nil {
		// if errors.Is(err, jwt.ErrTokenExpired) {
		// 	return "", errors.New("token is expired")
		// }
		return "", "", errors.New("token is expired")
	}

	if !token.Valid {
		return "", "", errors.New("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	if claims["type"] != "access" {
		return "", "", errors.New("invalid token type")
	}

	sub, _ := claims["sub"].(string)
	roleStr, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("invalid role claim")
	}

	role := user.Role(roleStr)

	return sub, role, nil
}

// Валидация Refresh Token
func (j *JwtUsecase) ValidateRefresh(tokenStr string) (string, user.Role, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secretRefresh, nil
	})

	if err != nil {
		// if errors.Is(err, jwt.ErrTokenExpired) {
		// 	return "", errors.New("token is expired")
		// }
		return "", "", errors.New("token is expired")
	}

	if !t.Valid {
		return "", "", err
	}

	claims := t.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	sub, _ := claims["sub"].(string)
	roleStr, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("invalid role claim")
	}

	role := user.Role(roleStr)
	return sub, role, nil
}

// Валидация Invite Token
func (j *JwtUsecase) ValidateInvite(tokenStr string) (string, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secretAccess, nil
	})
	if err != nil || !t.Valid {
		return "", err
	}

	claims := t.Claims.(jwt.MapClaims)

	if claims["type"] != "invite" {
		return "", errors.New("invalid token type")
	}

	return claims["email"].(string), nil
}
