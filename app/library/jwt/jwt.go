package jwt

import (
	"fmt"
	"goster/config"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID      string `json:"id"`
	IsAdmin bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func (c *Claims) GetUserId() (uuid.UUID, error) {
	return uuid.FromString(c.ID)
}

func GenerateToken(userID uuid.UUID, isAdmin bool, hours int) (string, error) {
	if hours <= 0 {
		hours = 4
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(hours))

	claims := Claims{
		ID:      userID.String(),
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "Goster",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(config.JWTsecret))
	if err != nil {
		return "", fmt.Errorf("Failed to sign token: %w", err)
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTsecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token or claims")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
