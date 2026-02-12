package utils

import (
	"robot_scheduler/internal/model/entity"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_Success(t *testing.T) {
	userID := uint(1)
	username := "testuser"
	role := entity.RoleUser
	secret := "test_secret_key"
	expireHours := 24

	token, err := GenerateToken(userID, username, role, secret, expireHours)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestValidateToken_Success(t *testing.T) {
	userID := uint(1)
	username := "testuser"
	role := entity.RoleUser
	secret := "test_secret_key"
	expireHours := 24

	// Generate token
	token, err := GenerateToken(userID, username, role, secret, expireHours)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate token
	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.UserName != username {
		t.Errorf("Expected UserName %s, got %s", username, claims.UserName)
	}

	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	userID := uint(1)
	username := "testuser"
	role := entity.RoleUser
	secret := "test_secret_key"

	// Create an expired token manually
	expirationTime := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	claims := &Claims{
		UserID:   userID,
		UserName: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	// Validate expired token
	_, err = ValidateToken(tokenString, secret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	userID := uint(1)
	username := "testuser"
	role := entity.RoleUser
	secret := "test_secret_key"
	wrongSecret := "wrong_secret_key"
	expireHours := 24

	// Generate token with one secret
	token, err := GenerateToken(userID, username, role, secret, expireHours)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	// Validate with different secret
	_, err = ValidateToken(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for invalid signature, got nil")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	secret := "test_secret_key"
	invalidToken := "invalid.token.string"

	_, err := ValidateToken(invalidToken, secret)
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	secret := "test_secret_key"

	_, err := ValidateToken("", secret)
	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}
}
