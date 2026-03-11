package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

const refreshTokenBytes = 32

func (s *Store) CreateRefreshToken(ctx context.Context, userID uuid.UUID, ttl time.Duration) (string, error) {
	token, hash, err := generateRefreshToken()
	if err != nil {
		return "", err
	}

	item := models.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(ttl),
	}

	if err := s.DB.WithContext(ctx).Create(&item).Error; err != nil {
		return "", err
	}
	return token, nil
}

func (s *Store) RotateRefreshToken(ctx context.Context, token string, ttl time.Duration) (uuid.UUID, string, error) {
	hash := hashRefreshToken(token)

	var userID uuid.UUID
	var newToken string
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.RefreshToken
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ? AND revoked_at IS NULL AND expires_at > now()", hash).
			First(&item).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return ErrInvalidRefreshToken
			}
			return err
		}

		now := time.Now()
		if err := tx.Model(&models.RefreshToken{}).
			Where("id = ?", item.ID).
			Updates(map[string]any{
				"revoked_at":   now,
				"last_used_at": now,
			}).Error; err != nil {
			return err
		}

		tokenValue, hashValue, err := generateRefreshToken()
		if err != nil {
			return err
		}

		next := models.RefreshToken{
			UserID:    item.UserID,
			TokenHash: hashValue,
			ExpiresAt: time.Now().Add(ttl),
		}
		if err := tx.Create(&next).Error; err != nil {
			return err
		}

		userID = item.UserID
		newToken = tokenValue
		return nil
	})
	if err != nil {
		return uuid.Nil, "", err
	}
	return userID, newToken, nil
}

func generateRefreshToken() (string, string, error) {
	buf := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	token := base64.RawURLEncoding.EncodeToString(buf)
	return token, hashRefreshToken(token), nil
}

func hashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
