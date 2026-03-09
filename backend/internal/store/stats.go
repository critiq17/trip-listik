package store

import (
	"context"

	"github.com/google/uuid"
)

type UserStatsSummary struct {
	TotalTrips       int64
	CountriesVisited int64
	CitiesVisited    int64
	TripsWithFriends int64
	SoloTrips        int64
}

func (s *Store) ComputeUserStats(ctx context.Context, userID uuid.UUID) (*UserStatsSummary, error) {
	var stats UserStatsSummary

	if err := s.DB.WithContext(ctx).Raw(`
		SELECT
			COUNT(DISTINCT t.id) AS total_trips
		FROM trips t
		JOIN trip_members tm ON tm.trip_id = t.id
		WHERE tm.user_id = ? AND t.deleted_at IS NULL
	`, userID).Scan(&stats).Error; err != nil {
		return nil, err
	}

	if err := s.DB.WithContext(ctx).Raw(`
		SELECT
			COUNT(DISTINCT t.country_code) AS countries_visited,
			COUNT(DISTINCT t.city) AS cities_visited
		FROM trips t
		JOIN trip_members tm ON tm.trip_id = t.id
		WHERE tm.user_id = ? AND t.status = 'completed' AND t.deleted_at IS NULL
	`, userID).Scan(&stats).Error; err != nil {
		return nil, err
	}

	if err := s.DB.WithContext(ctx).Raw(`
		SELECT
			COUNT(*) FILTER (WHERE member_count > 1) AS trips_with_friends,
			COUNT(*) FILTER (WHERE member_count = 1) AS solo_trips
		FROM (
			SELECT t.id, COUNT(tm.user_id) AS member_count
			FROM trips t
			JOIN trip_members tm ON tm.trip_id = t.id
			WHERE t.deleted_at IS NULL
			GROUP BY t.id
		) x
		JOIN trip_members tmu ON tmu.trip_id = x.id
		WHERE tmu.user_id = ?
	`, userID).Scan(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
