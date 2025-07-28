package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

type UserClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func (m *JWTManager) GenerateToken(userID string) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", err
	}

	return claims.UserID, nil
}
