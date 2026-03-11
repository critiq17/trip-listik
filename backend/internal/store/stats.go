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
		WITH user_trips AS (
			SELECT t.id, t.country_code, t.city, t.status
			FROM trips t
			JOIN trip_members tm ON tm.trip_id = t.id
			WHERE tm.user_id = ? AND t.deleted_at IS NULL
		),
		member_counts AS (
			SELECT tm.trip_id, COUNT(*) AS member_count
			FROM trip_members tm
			JOIN user_trips ut ON ut.id = tm.trip_id
			GROUP BY tm.trip_id
		)
		SELECT
			COUNT(*) AS total_trips,
			COUNT(DISTINCT CASE WHEN ut.status = 'completed' THEN ut.country_code END) AS countries_visited,
			COUNT(DISTINCT CASE WHEN ut.status = 'completed' THEN ut.city END) AS cities_visited,
			COUNT(*) FILTER (WHERE mc.member_count > 1) AS trips_with_friends,
			COUNT(*) FILTER (WHERE mc.member_count = 1) AS solo_trips
		FROM user_trips ut
		LEFT JOIN member_counts mc ON mc.trip_id = ut.id
	`, userID).Scan(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
