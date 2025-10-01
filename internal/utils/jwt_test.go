package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	tokenString, err := GenerateAccessToken(123)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatalf("Error converting claims")
	}

	if claims.UserID != 123 {
		t.Errorf("Expected UserID 123, received %v", claims.UserID)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Errorf("Token has already expired")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenString, err := GenerateRefreshToken(1234)
	if err != nil {
		t.Fatalf("Error generating refresh token")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		t.Fatalf("Error converting claim")
	}

	if claims.UserID != 1234 {
		t.Errorf("Expected UserID 1234, received %v", claims.UserID)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Errorf("Token has already expired")
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "mySecurePassword"

	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	assert.True(t, CheckPasswordHash(password, hash))

	assert.False(t, CheckPasswordHash("wrongPassword", hash))

	originalCost := BcryptCost
	BcryptCost = 100
	defer func() { BcryptCost = originalCost }()

	_, err = HashPassword("any")
	assert.Error(t, err)
}
