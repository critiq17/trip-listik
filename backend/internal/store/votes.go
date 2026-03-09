package store

import (
	"context"

	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type VoteSummary struct {
	Average float64
	Count   int64
}

func (s *Store) UpsertVote(ctx context.Context, tripID, userID uuid.UUID, vote int16) error {
	v := models.TripVote{TripID: tripID, UserID: userID, Vote: vote}
	return s.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "trip_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"vote", "updated_at"}),
		}).
		Create(&v).Error
}

func (s *Store) GetVoteSummary(ctx context.Context, tripID uuid.UUID) (*VoteSummary, error) {
	var summary VoteSummary
	if err := s.DB.WithContext(ctx).
		Model(&models.TripVote{}).
		Select("COALESCE(AVG(vote),0) as average, COUNT(*) as count").
		Where("trip_id = ?", tripID).
		Scan(&summary).Error; err != nil {
		return nil, err
	}
	return &summary, nil
}
