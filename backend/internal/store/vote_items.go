package store

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store/models"
)

// ── VoteItem store methods ────────────────────────────────────────────────────

type VoteItemView struct {
	ID          string    `json:"id"`
	TripID      string    `json:"trip_id"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	AddedBy     string    `json:"added_by"`
	VoteCount   int64     `json:"vote_count"`
	HasVoted    bool      `json:"has_voted"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Store) CreateVoteItem(ctx context.Context, item *models.VoteItem) error {
	return s.DB.WithContext(ctx).Create(item).Error
}

func (s *Store) GetVoteItemByID(ctx context.Context, id string) (*models.VoteItem, error) {
	var item models.VoteItem
	if err := s.DB.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Store) ListVoteItems(ctx context.Context, tripID string, viewerUserID string) ([]VoteItemView, error) {
	var items []VoteItemView
	err := s.DB.WithContext(ctx).
		Table("vote_items vi").
		Select(`vi.id, vi.trip_id, vi.category, vi.title, vi.description, vi.image_url, vi.added_by, vi.created_at,
			COUNT(v.item_id) AS vote_count,
			BOOL_OR(v.user_id = ?) AS has_voted`, viewerUserID).
		Joins("LEFT JOIN vote_item_votes v ON v.item_id = vi.id").
		Where("vi.trip_id = ?", tripID).
		Group("vi.id").
		Order("vi.created_at ASC").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Store) DeleteVoteItem(ctx context.Context, id string) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.VoteItem{}).Error
}

// ── VoteItemVote ─────────────────────────────────────────────────────────────

func (s *Store) AddVoteItemVote(ctx context.Context, itemID, userID string) error {
	return s.DB.WithContext(ctx).
		Exec("INSERT INTO vote_item_votes (item_id, user_id, created_at) VALUES (?, ?, NOW()) ON CONFLICT DO NOTHING",
			itemID, userID).Error
}

func (s *Store) RemoveVoteItemVote(ctx context.Context, itemID, userID string) error {
	return s.DB.WithContext(ctx).
		Exec("DELETE FROM vote_item_votes WHERE item_id = ? AND user_id = ?", itemID, userID).Error
}

func (s *Store) GetVoteItemCount(ctx context.Context, itemID string) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).
		Model(&models.VoteItemVote{}).
		Where("item_id = ?", itemID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
