package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

// CustomClaims can hold any extra fields you want.
type CustomClaims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, tokenType string) (string, error) {
	now := time.Now()

	claims := CustomClaims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now), // token not valid before now
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ParseAndVerifyJWT(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the algorithm is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
