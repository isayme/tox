package util

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

func ValidateJwtToken(tokenString string, key []byte) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("not valid token")
	}

	return nil
}

func GenerateJwtToken(key []byte) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: NowInMills() + 1500,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}
