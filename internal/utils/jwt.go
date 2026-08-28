package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManagerInterface interface {
	GenerateTokens(userID int) (string, error)
	ParseToken(tokenString string) (int, error)
}

type JWTManager struct {
	secret string
}

type Claims struct {
	UserID int

	jwt.RegisteredClaims
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret,
	}
}

func (j *JWTManager) GenerateTokens(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (j *JWTManager) ParseToken(tokenString string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return 0, err
	}

	if !parsedToken.Valid {
		return 0, errors.New("token is invalid")
	}

	claims, ok := parsedToken.Claims.(*Claims)

	if !ok {
		return 0, errors.New("invalid claims")
	}

	return claims.UserID, nil
}
