package store

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type UserSummary struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PhotoURL  string    `json:"photo_url"`
}

func (s *Store) SearchUsers(ctx context.Context, query string, limit int) ([]UserSummary, error) {
	var users []UserSummary
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return users, nil
	}
	pattern := "%" + escapeLike(trimmed) + "%"
	db := s.DB.WithContext(ctx).
		Table("users").
		Select("id, username, first_name, last_name, photo_url").
		Where("username ILIKE ? ESCAPE '\\' OR first_name ILIKE ? ESCAPE '\\' OR last_name ILIKE ? ESCAPE '\\'", pattern, pattern, pattern).
		Order("username ASC")
	if limit > 0 {
		db = db.Limit(limit)
	}
	if err := db.Scan(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
