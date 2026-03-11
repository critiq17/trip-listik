package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FeedHandler struct {
	Store *store.Store
}

type TripFeedItem struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	CoverPhotoURL string     `json:"cover_photo_url"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	CountryCode   string     `json:"country_code"`
	City          string     `json:"city"`
	MemberCount   int64      `json:"member_count"`
	VoteCount     int64      `json:"vote_count"`
	VoteAverage   float64    `json:"vote_average"`
}

func (h *FeedHandler) Feed(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	filter := c.Query("filter", "all")
	country := c.Query("country", "")
	cursorStr := c.Query("cursor", "")
	var cursor time.Time
	if cursorStr != "" {
		parsed, err := time.Parse(time.RFC3339Nano, cursorStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid cursor")
		}
		cursor = parsed
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	var trips []models.Trip
	var err error

	switch filter {
	case "friends":
		userID, ok := middleware.GetUserID(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		trips, err = h.Store.ListFriendTrips(ctx, userID, limit, cursor)
	case "popular":
		trips, err = h.Store.ListPopularTrips(ctx, limit, cursor)
	case "nearby":
		if country == "" {
			return fiber.NewError(fiber.StatusBadRequest, "country is required for nearby filter")
		}
		trips, err = h.Store.ListNearbyTrips(ctx, country, limit, cursor)
	default:
		trips, err = h.Store.ListPublicTrips(ctx, limit, cursor)
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load feed")
	}

	tripIDs := make([]uuid.UUID, 0, len(trips))
	for _, t := range trips {
		tripIDs = append(tripIDs, t.ID)
	}

	memberCounts, err := h.Store.GetTripMemberCounts(ctx, tripIDs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load member counts")
	}
	voteStats, err := h.Store.GetTripVoteStats(ctx, tripIDs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load vote stats")
	}

	items := make([]TripFeedItem, 0, len(trips))
	for _, t := range trips {
		vs := voteStats[t.ID]
		items = append(items, TripFeedItem{
			ID:            t.ID,
			Title:         t.Title,
			CoverPhotoURL: t.CoverPhotoURL,
			StartDate:     t.StartDate,
			EndDate:       t.EndDate,
			CountryCode:   t.CountryCode,
			City:          t.City,
			MemberCount:   memberCounts[t.ID],
			VoteCount:     vs.Count,
			VoteAverage:   vs.Average,
		})
	}

	nextCursor := ""
	if len(trips) == limit {
		last := trips[len(trips)-1]
		nextCursor = last.CreatedAt.UTC().Format(time.RFC3339Nano)
	}

	return c.JSON(fiber.Map{
		"items":       items,
		"next_cursor": nextCursor,
	})
}

func (h *FeedHandler) Explore(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	query := c.Query("q", "")
	country := c.Query("country", "")
	cursorStr := c.Query("cursor", "")
	var cursor time.Time
	if cursorStr != "" {
		parsed, err := time.Parse(time.RFC3339Nano, cursorStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid cursor")
		}
		cursor = parsed
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	trips, err := h.Store.SearchPublicTrips(ctx, query, country, limit, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load explore feed")
	}

	tripIDs := make([]uuid.UUID, 0, len(trips))
	for _, t := range trips {
		tripIDs = append(tripIDs, t.ID)
	}

	memberCounts, err := h.Store.GetTripMemberCounts(ctx, tripIDs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load member counts")
	}
	voteStats, err := h.Store.GetTripVoteStats(ctx, tripIDs)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load vote stats")
	}

	items := make([]TripFeedItem, 0, len(trips))
	for _, t := range trips {
		vs := voteStats[t.ID]
		items = append(items, TripFeedItem{
			ID:            t.ID,
			Title:         t.Title,
			CoverPhotoURL: t.CoverPhotoURL,
			StartDate:     t.StartDate,
			EndDate:       t.EndDate,
			CountryCode:   t.CountryCode,
			City:          t.City,
			MemberCount:   memberCounts[t.ID],
			VoteCount:     vs.Count,
			VoteAverage:   vs.Average,
		})
	}

	nextCursor := ""
	if len(trips) == limit {
		last := trips[len(trips)-1]
		nextCursor = last.CreatedAt.UTC().Format(time.RFC3339Nano)
	}

	return c.JSON(fiber.Map{
		"items":       items,
		"next_cursor": nextCursor,
	})
}
