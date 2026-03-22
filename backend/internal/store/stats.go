package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserStatsSummary struct {
	TotalTrips       int64 `json:"total_trips"`
	CountriesVisited int64 `json:"countries_visited"`
	CitiesVisited    int64 `json:"cities_visited"`
	TripsWithFriends int64 `json:"trips_with_friends"`
	SoloTrips        int64 `json:"solo_trips"`
}

type CountryVisit struct {
	Code       string    `json:"code"`
	VisitCount int64     `json:"visit_count"`
	LastVisit  time.Time `json:"last_visit"`
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

func (s *Store) GetUserCountryVisits(ctx context.Context, userID uuid.UUID) ([]CountryVisit, error) {
	var rows []CountryVisit
	if err := s.DB.WithContext(ctx).Raw(`
		SELECT
			t.country_code AS code,
			COUNT(*) AS visit_count,
			MAX(COALESCE(t.end_date, t.start_date, t.created_at)) AS last_visit
		FROM trips t
		JOIN trip_members tm ON tm.trip_id = t.id
		WHERE tm.user_id = ? AND t.country_code IS NOT NULL AND t.country_code <> '' AND t.deleted_at IS NULL
		GROUP BY t.country_code
		ORDER BY visit_count DESC, last_visit DESC
	`, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
