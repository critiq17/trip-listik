package handlers

// Integration tests for VoteItemsHandler.Vote
//
// Set TEST_DATABASE_URL to run these tests, e.g.:
//
//	TEST_DATABASE_URL="host=localhost user=myuser password=1234 dbname=trip_listik_test sslmode=disable"
//
// Without that env var every test is skipped automatically.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/critiq17/tripListik/internal/db"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ── helpers ──────────────────────────────────────────────────────────────────

// setupVoteStore connects to the test database and skips if the env var is absent.
func setupVoteStore(t *testing.T) *store.Store {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set — skipping integration tests")
	}
	database, _, err := db.Connect(dsn)
	if err != nil {
		t.Fatalf("db.Connect: %v", err)
	}
	return store.New(database)
}

// seedVoteUser inserts a minimal user row and returns it.
func seedVoteUser(t *testing.T, db *gorm.DB) *models.User {
	t.Helper()
	u := &models.User{
		TelegramID: int64(uuid.New().ID()),
		Username:   uuid.New().String()[:10],
		FirstName:  "Test",
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seedVoteUser: %v", err)
	}
	return u
}

// seedPrivateTrip creates a private trip owned by ownerID and adds ownerID as a member.
func seedPrivateTrip(t *testing.T, s *store.Store, ownerID uuid.UUID) *models.Trip {
	t.Helper()
	trip := &models.Trip{
		OwnerID:    ownerID,
		Title:      "Vote Test Trip",
		Visibility: "private",
		Status:     "active",
	}
	if err := s.DB.Create(trip).Error; err != nil {
		t.Fatalf("seedPrivateTrip: create: %v", err)
	}
	if err := s.AddTripMember(context.Background(), trip.ID, ownerID, "owner"); err != nil {
		t.Fatalf("seedPrivateTrip: add member: %v", err)
	}
	return trip
}

// seedVoteItem creates a vote_items row for the given trip and returns its ID.
func seedVoteItem(t *testing.T, s *store.Store, tripID uuid.UUID, addedBy uuid.UUID) string {
	t.Helper()
	item := &models.VoteItem{
		TripID:   tripID.String(),
		Category: "ACTIVITY",
		Title:    "Rafting",
		AddedBy:  addedBy.String(),
	}
	if err := s.CreateVoteItem(context.Background(), item); err != nil {
		t.Fatalf("seedVoteItem: %v", err)
	}
	return item.ID
}

// buildVoteApp constructs a minimal Fiber app that routes
//
//	POST /v1/trips/:id/vote-items/:itemId/vote
//
// with a fake auth middleware that injects callerID into fiber locals so that
// GetUserID() works without requiring a real JWT.
func buildVoteApp(s *store.Store, callerID uuid.UUID) *fiber.App {
	app := fiber.New(fiber.Config{
		// Suppress the default error handler's HTML output in tests.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Fake auth: inject the caller's UUID into locals, mimicking RequireAuth.
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", callerID)
		return c.Next()
	})

	h := &VoteItemsHandler{Store: s, Hub: nil}
	app.Post("/v1/trips/:id/vote-items/:itemId/vote", h.Vote)

	return app
}

// doVoteRequest fires a POST to the vote endpoint and returns the response.
func doVoteRequest(t *testing.T, app *fiber.App, tripID uuid.UUID, itemID string) *http.Response {
	t.Helper()
	url := fmt.Sprintf("/v1/trips/%s/vote-items/%s/vote", tripID, itemID)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	resp, err := app.Test(req, 5000 /* ms */)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp
}

// ── TestVoteItemsHandler_Vote_ForbiddenNonMember ──────────────────────────────

// TestVoteItemsHandler_Vote_ForbiddenNonMember verifies that a user who is not
// a member of a private trip receives HTTP 403 when attempting to vote.
func TestVoteItemsHandler_Vote_ForbiddenNonMember(t *testing.T) {
	s := setupVoteStore(t)

	// Create the trip owner and the trip itself.
	owner := seedVoteUser(t, s.DB)
	trip := seedPrivateTrip(t, s, owner.ID)
	itemID := seedVoteItem(t, s, trip.ID, owner.ID)

	// The caller is a completely different user who has never been added as a
	// member of this trip.
	nonMember := seedVoteUser(t, s.DB)

	app := buildVoteApp(s, nonMember.ID)
	resp := doVoteRequest(t, app, trip.ID, itemID)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status 403, got %d; body: %s", resp.StatusCode, body)
	}
}

// ── TestVoteItemsHandler_Vote_IdempotentDoubleVote ────────────────────────────

// TestVoteItemsHandler_Vote_IdempotentDoubleVote verifies that voting twice with
// the same user on the same item is idempotent: both calls return HTTP 200 and
// the database records exactly one vote (ON CONFLICT DO NOTHING semantics).
func TestVoteItemsHandler_Vote_IdempotentDoubleVote(t *testing.T) {
	s := setupVoteStore(t)

	owner := seedVoteUser(t, s.DB)
	trip := seedPrivateTrip(t, s, owner.ID)
	itemID := seedVoteItem(t, s, trip.ID, owner.ID)

	// Clean up the vote rows created by this test so parallel or repeated runs
	// don't accumulate stale data.
	t.Cleanup(func() {
		s.DB.Exec("DELETE FROM vote_item_votes WHERE item_id = ?", itemID)
		s.DB.Exec("DELETE FROM vote_items WHERE id = ?", itemID)
		s.DB.Exec("DELETE FROM trip_members WHERE trip_id = ?", trip.ID)
		s.DB.Exec("DELETE FROM trips WHERE id = ?", trip.ID)
		s.DB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, owner.ID)
	})

	app := buildVoteApp(s, owner.ID)

	// First vote — must succeed.
	resp1 := doVoteRequest(t, app, trip.ID, itemID)
	defer resp1.Body.Close()
	if resp1.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp1.Body)
		t.Fatalf("first vote: expected 200, got %d; body: %s", resp1.StatusCode, body)
	}

	// Second vote with the same user — must also succeed (no error from DB).
	resp2 := doVoteRequest(t, app, trip.ID, itemID)
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp2.Body)
		t.Fatalf("second vote: expected 200, got %d; body: %s", resp2.StatusCode, body)
	}

	// Parse the response body of the second call to inspect vote_count.
	var result struct {
		VoteCount int64 `json:"vote_count"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.VoteCount != 1 {
		t.Errorf("expected vote_count=1 after two identical votes, got %d", result.VoteCount)
	}

	// Double-check directly via the store to be sure the HTTP layer and the DB
	// agree.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	count, err := s.GetVoteItemCount(ctx, itemID)
	if err != nil {
		t.Fatalf("GetVoteItemCount: %v", err)
	}
	if count != 1 {
		t.Errorf("DB vote_count = %d, want 1", count)
	}
}
