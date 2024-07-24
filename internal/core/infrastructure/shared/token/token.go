package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/config"
)

type token struct {
	signingKey string
}

func NewToken(config *config.Config) TokenService {
	return &token{
		signingKey: config.TokenConfig.SignKey,
	}
}

func (t *token) CreateJwtToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userID,
		"role": "patient",
		"exp":  time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(t.signingKey))
}
