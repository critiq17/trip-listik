package trips_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/critiq17/tripListik/internal/db"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/trips"
	"github.com/google/uuid"
)

// setupSvc initialises a real test database.
// Set TEST_DATABASE_URL to run these tests, e.g.:
//   TEST_DATABASE_URL="host=localhost user=myuser password=1234 dbname=trip_listik_test sslmode=disable"
func setupSvc(t *testing.T) *trips.Service {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set — skipping integration tests")
	}
	database, _, err := db.Connect(dsn)
	if err != nil {
		t.Fatalf("db.Connect: %v", err)
	}
	return trips.NewService(store.New(database))
}

// ── CreateTrip ──────────────────────────────────────────────────────────────

func TestCreateTrip_Valid(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{
		Title: "Paris 2026",
		City:  "Paris",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trip.ID == uuid.Nil {
		t.Error("expected non-nil trip ID")
	}
	if trip.Title != "Paris 2026" {
		t.Errorf("title = %q, want %q", trip.Title, "Paris 2026")
	}
	if trip.OwnerID != ownerID {
		t.Error("owner_id mismatch")
	}
}

func TestCreateTrip_EmptyTitle(t *testing.T) {
	svc := setupSvc(t)

	_, err := svc.CreateTrip(context.Background(), uuid.New(), trips.CreateTripInput{
		Title: "   ",
	})
	if err == nil {
		t.Fatal("expected validation error for blank title")
	}
	if !strings.Contains(err.Error(), "title") {
		t.Errorf("error should mention 'title', got: %v", err)
	}
}

func TestCreateTrip_TitleTooLong(t *testing.T) {
	svc := setupSvc(t)

	_, err := svc.CreateTrip(context.Background(), uuid.New(), trips.CreateTripInput{
		Title: strings.Repeat("a", 101),
	})
	if err == nil {
		t.Fatal("expected validation error for title > 100 chars")
	}
}

func TestCreateTrip_DescriptionTooLong(t *testing.T) {
	svc := setupSvc(t)

	_, err := svc.CreateTrip(context.Background(), uuid.New(), trips.CreateTripInput{
		Title:       "Valid",
		Description: strings.Repeat("x", 2001),
	})
	if err == nil {
		t.Fatal("expected validation error for description > 2000 chars")
	}
}

func TestCreateTrip_EndBeforeStart(t *testing.T) {
	svc := setupSvc(t)

	_, err := svc.CreateTrip(context.Background(), uuid.New(), trips.CreateTripInput{
		Title:     "Bad Dates",
		StartDate: "2026-06-10",
		EndDate:   "2026-06-05",
	})
	if err == nil {
		t.Fatal("expected validation error when end_date is before start_date")
	}
	if !strings.Contains(err.Error(), "end_date") {
		t.Errorf("error should mention end_date, got: %v", err)
	}
}

func TestCreateTrip_SameDayStartEnd(t *testing.T) {
	svc := setupSvc(t)

	_, err := svc.CreateTrip(context.Background(), uuid.New(), trips.CreateTripInput{
		Title:     "Day Trip",
		StartDate: "2026-06-10",
		EndDate:   "2026-06-10",
	})
	if err != nil {
		t.Fatalf("same-day start/end should be valid: %v", err)
	}
}

// ── DeleteTrip ──────────────────────────────────────────────────────────────

func TestDeleteTrip_NotOwner(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{Title: "To delete"})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	err = svc.DeleteTrip(context.Background(), uuid.New(), trip.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden when non-owner deletes")
	}
	if err != trips.ErrForbidden {
		t.Errorf("want ErrForbidden, got: %v", err)
	}
}

func TestDeleteTrip_NotFound(t *testing.T) {
	svc := setupSvc(t)

	err := svc.DeleteTrip(context.Background(), uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected ErrNotFound for unknown trip")
	}
}

func TestDeleteTrip_Owner(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{Title: "Delete me"})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := svc.DeleteTrip(context.Background(), ownerID, trip.ID); err != nil {
		t.Fatalf("owner should be able to delete: %v", err)
	}

	_, err = svc.GetTrip(context.Background(), trip.ID, &ownerID)
	if err == nil {
		t.Fatal("deleted trip should not be found")
	}
}

// ── UpdateTrip ──────────────────────────────────────────────────────────────

func TestUpdateTrip_Forbidden(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, _ := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{Title: "Original"})

	newTitle := "Hacked"
	_, err := svc.UpdateTrip(context.Background(), uuid.New(), trip.ID, trips.UpdateTripInput{Title: &newTitle})
	if err == nil {
		t.Fatal("expected ErrForbidden when non-owner updates")
	}
}

func TestUpdateTrip_EndBeforeExistingStart(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, _ := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{
		Title:     "Dated trip",
		StartDate: "2026-05-01",
	})

	badEnd := "2026-04-01"
	_, err := svc.UpdateTrip(context.Background(), ownerID, trip.ID, trips.UpdateTripInput{EndDate: &badEnd})
	if err == nil {
		t.Fatal("expected validation error: end_date before existing start_date")
	}
}
