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

// GetUserFriends returns distinct users who share at least one non-deleted trip with the given user.
func (s *Store) GetUserFriends(ctx context.Context, userID uuid.UUID) ([]UserSummary, error) {
	var friends []UserSummary
	err := s.DB.WithContext(ctx).
		Table("users u").
		Select("DISTINCT u.id, u.username, u.first_name, u.last_name, u.photo_url").
		Joins("JOIN trip_members tm1 ON tm1.user_id = u.id").
		Joins("JOIN trip_members tm2 ON tm2.trip_id = tm1.trip_id AND tm2.user_id = ?", userID).
		Joins("JOIN trips t ON t.id = tm1.trip_id AND t.deleted_at IS NULL").
		Where("u.id != ?", userID).
		Order("u.username ASC").
		Limit(100).
		Scan(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
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
