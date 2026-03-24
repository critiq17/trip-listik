package trips_test

import (
	"context"
	"errors"
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
//
//	TEST_DATABASE_URL="host=localhost user=myuser password=1234 dbname=trip_listik_test sslmode=disable"
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

// nilSvc returns a Service with a nil store. Safe only for test cases whose
// input triggers a validation error before any store method is called.
func nilSvc() *trips.Service {
	return trips.NewService(nil)
}

// ── CreateTrip unit tests (no DB required) ─────────────────────────────────

func TestCreateTrip_Validation(t *testing.T) {
	t.Parallel()

	longTitle := strings.Repeat("a", 101)
	longDesc := strings.Repeat("x", 2001)

	tests := []struct {
		name    string
		input   trips.CreateTripInput
		wantErr error
	}{
		{
			name:    "empty title",
			input:   trips.CreateTripInput{Title: ""},
			wantErr: trips.ErrValidation,
		},
		{
			name:    "whitespace-only title",
			input:   trips.CreateTripInput{Title: "   "},
			wantErr: trips.ErrValidation,
		},
		{
			name:    "title too long",
			input:   trips.CreateTripInput{Title: longTitle},
			wantErr: trips.ErrValidation,
		},
		{
			name:    "description too long",
			input:   trips.CreateTripInput{Title: "Valid", Description: longDesc},
			wantErr: trips.ErrValidation,
		},
		{
			name: "end date before start date",
			input: trips.CreateTripInput{
				Title:     "Bad Dates",
				StartDate: "2026-06-10",
				EndDate:   "2026-06-05",
			},
			wantErr: trips.ErrValidation,
		},
		{
			name: "invalid start date format",
			input: trips.CreateTripInput{
				Title:     "Bad Format",
				StartDate: "10-06-2026",
			},
			wantErr: trips.ErrValidation,
		},
		{
			name: "invalid end date format",
			input: trips.CreateTripInput{
				Title:   "Bad End Format",
				EndDate: "not-a-date",
			},
			wantErr: trips.ErrValidation,
		},
	}

	svc := nilSvc()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := svc.CreateTrip(context.Background(), uuid.New(), tc.input)
			if err == nil {
				t.Fatalf("expected error wrapping %v, got nil", tc.wantErr)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("errors.Is(%v, %v) = false; full error: %v", err, tc.wantErr, err)
			}
		})
	}
}

// ── UpdateTrip unit tests (no DB required) ─────────────────────────────────

// TestUpdateTrip_Validation tests input validation cases inside UpdateTrip
// that do NOT require a real database because the store is called first to
// fetch the trip. These cases therefore require an integration DB — they are
// grouped here with the integration tests below. For purely pre-fetch
// validation we rely on the integration path.
//
// The table below covers cases that ARE reachable without a DB: the service
// validates field values AFTER fetching the trip from the store.  Since we
// cannot stub the store without a mock library, the validation-only UpdateTrip
// cases are exercised through the integration helpers below.
// This comment is intentionally kept to document the design decision.

// ── CreateTrip integration tests ────────────────────────────────────────────

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
	if !errors.Is(err, trips.ErrValidation) {
		t.Errorf("want ErrValidation, got: %v", err)
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
	if !errors.Is(err, trips.ErrValidation) {
		t.Errorf("want ErrValidation, got: %v", err)
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
	if !errors.Is(err, trips.ErrValidation) {
		t.Errorf("want ErrValidation, got: %v", err)
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
	if !errors.Is(err, trips.ErrValidation) {
		t.Errorf("want ErrValidation, got: %v", err)
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

func TestCreateTrip_DefaultVisibilityAndStatus(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{
		Title: "Defaults Check",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trip.Visibility != "public" {
		t.Errorf("default visibility = %q, want %q", trip.Visibility, "public")
	}
	if trip.Status != "draft" {
		t.Errorf("default status = %q, want %q", trip.Status, "draft")
	}
}

func TestCreateTrip_ExplicitVisibilityAndStatus(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{
		Title:      "Private Trip",
		Visibility: "private",
		Status:     "active",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if trip.Visibility != "private" {
		t.Errorf("visibility = %q, want %q", trip.Visibility, "private")
	}
	if trip.Status != "active" {
		t.Errorf("status = %q, want %q", trip.Status, "active")
	}
}

// ── DeleteTrip integration tests ────────────────────────────────────────────

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
	if !errors.Is(err, trips.ErrForbidden) {
		t.Errorf("want ErrForbidden, got: %v", err)
	}
}

func TestDeleteTrip_NotFound(t *testing.T) {
	svc := setupSvc(t)

	err := svc.DeleteTrip(context.Background(), uuid.New(), uuid.New())
	if err == nil {
		t.Fatal("expected ErrNotFound for unknown trip")
	}
	if !errors.Is(err, trips.ErrNotFound) {
		t.Errorf("want ErrNotFound, got: %v", err)
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

// ── UpdateTrip integration tests ────────────────────────────────────────────

func TestUpdateTrip_Validation(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	base, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{
		Title:     "Base Trip",
		StartDate: "2026-05-01",
	})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	emptyTitle := ""
	longTitle := strings.Repeat("a", 101)
	longDesc := strings.Repeat("x", 2001)
	badEnd := "2026-04-01"
	badDateFmt := "01/05/2026"

	tests := []struct {
		name    string
		input   trips.UpdateTripInput
		wantErr error
	}{
		{
			name:    "empty title",
			input:   trips.UpdateTripInput{Title: &emptyTitle},
			wantErr: trips.ErrValidation,
		},
		{
			name:    "title too long",
			input:   trips.UpdateTripInput{Title: &longTitle},
			wantErr: trips.ErrValidation,
		},
		{
			name:    "description too long",
			input:   trips.UpdateTripInput{Description: &longDesc},
			wantErr: trips.ErrValidation,
		},
		{
			name:    "end date before existing start date",
			input:   trips.UpdateTripInput{EndDate: &badEnd},
			wantErr: trips.ErrValidation,
		},
		{
			name:    "invalid end date format",
			input:   trips.UpdateTripInput{EndDate: &badDateFmt},
			wantErr: trips.ErrValidation,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.UpdateTrip(context.Background(), ownerID, base.ID, tc.input)
			if err == nil {
				t.Fatalf("expected error wrapping %v, got nil", tc.wantErr)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("errors.Is(%v, %v) = false; full error: %v", err, tc.wantErr, err)
			}
		})
	}
}

func TestUpdateTrip_Forbidden(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{Title: "Original"})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	newTitle := "Hacked"
	_, err = svc.UpdateTrip(context.Background(), uuid.New(), trip.ID, trips.UpdateTripInput{Title: &newTitle})
	if err == nil {
		t.Fatal("expected ErrForbidden when non-owner updates")
	}
	if !errors.Is(err, trips.ErrForbidden) {
		t.Errorf("want ErrForbidden, got: %v", err)
	}
}

func TestUpdateTrip_NotFound(t *testing.T) {
	svc := setupSvc(t)

	newTitle := "Ghost"
	_, err := svc.UpdateTrip(context.Background(), uuid.New(), uuid.New(), trips.UpdateTripInput{Title: &newTitle})
	if err == nil {
		t.Fatal("expected ErrNotFound for unknown trip")
	}
	if !errors.Is(err, trips.ErrNotFound) {
		t.Errorf("want ErrNotFound, got: %v", err)
	}
}

func TestUpdateTrip_ValidByOwner(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{Title: "Original Title"})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	updated := "Updated Title"
	city := "Rome"
	result, err := svc.UpdateTrip(context.Background(), ownerID, trip.ID, trips.UpdateTripInput{
		Title: &updated,
		City:  &city,
	})
	if err != nil {
		t.Fatalf("owner update should succeed: %v", err)
	}
	if result.Title != updated {
		t.Errorf("title = %q, want %q", result.Title, updated)
	}
	if result.City != city {
		t.Errorf("city = %q, want %q", result.City, city)
	}
}

func TestUpdateTrip_EndBeforeExistingStart(t *testing.T) {
	svc := setupSvc(t)
	ownerID := uuid.New()

	trip, err := svc.CreateTrip(context.Background(), ownerID, trips.CreateTripInput{
		Title:     "Dated trip",
		StartDate: "2026-05-01",
	})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	badEnd := "2026-04-01"
	_, err = svc.UpdateTrip(context.Background(), ownerID, trip.ID, trips.UpdateTripInput{EndDate: &badEnd})
	if err == nil {
		t.Fatal("expected validation error: end_date before existing start_date")
	}
	if !errors.Is(err, trips.ErrValidation) {
		t.Errorf("want ErrValidation, got: %v", err)
	}
}
