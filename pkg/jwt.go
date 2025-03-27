package pkg

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func ParseToken(tokenString string, secret []byte) (string, error) {
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if !t.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	id, _ := claims["sub"].(string)
	if id == "" {
		return "", errors.New("invalid subject")
	}
	return id, nil
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Fields(authHeader)
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1], nil
	}

	if len(parts) == 1 {
		return parts[0], nil
	}

	return "", errors.New("invalid token format")
}
