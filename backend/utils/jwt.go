package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	UserID string
	Email  string
}

// GenerateToken generates an access token (valid for 1 hour) and a refresh token (valid for 7 days)
func GenerateToken(payload TokenPayload) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")

	// Access Token
	atClaims := jwt.MapClaims{
		"user_id": payload.UserID,
		"email":   payload.Email,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	rtClaims := jwt.MapClaims{
		"user_id": payload.UserID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// VerifyToken verifies a JWT token and returns the claims
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
