package store

import (
	"context"

	"github.com/google/uuid"
)

type TripMemberCount struct {
	TripID uuid.UUID
	Count  int64
}

type TripVoteStats struct {
	TripID  uuid.UUID
	Count   int64
	Average float64
}

type TripPhotoCount struct {
	TripID uuid.UUID
	Count  int64
}

type TripAggregates struct {
	MemberCount  int64
	VoteCount    int64
	VoteAverage  float64
	CommentCount int64
	PhotoCount   int64
}

func (s *Store) GetTripMemberCounts(ctx context.Context, tripIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	if len(tripIDs) == 0 {
		return map[uuid.UUID]int64{}, nil
	}
	var rows []TripMemberCount
	if err := s.DB.WithContext(ctx).
		Table("trip_members").
		Select("trip_id, COUNT(*) as count").
		Where("trip_id IN ?", tripIDs).
		Group("trip_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uuid.UUID]int64, len(rows))
	for _, r := range rows {
		out[r.TripID] = r.Count
	}
	return out, nil
}

func (s *Store) GetTripVoteStats(ctx context.Context, tripIDs []uuid.UUID) (map[uuid.UUID]TripVoteStats, error) {
	if len(tripIDs) == 0 {
		return map[uuid.UUID]TripVoteStats{}, nil
	}
	var rows []TripVoteStats
	if err := s.DB.WithContext(ctx).
		Table("trip_votes").
		Select("trip_id, COUNT(*) as count, COALESCE(AVG(vote),0) as average").
		Where("trip_id IN ?", tripIDs).
		Group("trip_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uuid.UUID]TripVoteStats, len(rows))
	for _, r := range rows {
		out[r.TripID] = r
	}
	return out, nil
}

func (s *Store) GetTripPhotoCounts(ctx context.Context, tripIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	if len(tripIDs) == 0 {
		return map[uuid.UUID]int64{}, nil
	}
	var rows []TripPhotoCount
	if err := s.DB.WithContext(ctx).
		Table("trip_photos").
		Select("trip_id, COUNT(*) as count").
		Where("trip_id IN ?", tripIDs).
		Group("trip_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[uuid.UUID]int64, len(rows))
	for _, r := range rows {
		out[r.TripID] = r.Count
	}
	return out, nil
}

func (s *Store) GetTripAggregates(ctx context.Context, tripID uuid.UUID) (*TripAggregates, error) {
	var agg TripAggregates
	if err := s.DB.WithContext(ctx).Raw(`
		SELECT
			(SELECT COUNT(*) FROM trip_members WHERE trip_id = ?) AS member_count,
			(SELECT COUNT(*) FROM trip_votes WHERE trip_id = ?) AS vote_count,
			(SELECT COALESCE(AVG(vote), 0) FROM trip_votes WHERE trip_id = ?) AS vote_average,
			(SELECT COUNT(*) FROM trip_comments WHERE trip_id = ?) AS comment_count,
			(SELECT COUNT(*) FROM trip_photos WHERE trip_id = ?) AS photo_count
	`, tripID, tripID, tripID, tripID, tripID).Scan(&agg).Error; err != nil {
		return nil, err
	}
	return &agg, nil
}
