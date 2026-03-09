package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestNewToken(t *testing.T) {
	userID := uuid.New()
	token, err := NewToken(userID, 12345, "secret", time.Minute)
	if err != nil {
		t.Fatalf("expected token, got error: %v", err)
	}

	claims := &Claims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		t.Fatalf("expected valid token, got error: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, claims.UserID)
	}
	if claims.TelegramID != 12345 {
		t.Fatalf("expected telegram id 12345, got %d", claims.TelegramID)
	}
}
