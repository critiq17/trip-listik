package invites_test

import (
	"context"
	"os"
	"testing"

	"github.com/critiq17/tripListik/internal/db"
	"github.com/critiq17/tripListik/internal/invites"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/critiq17/tripListik/internal/trips"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type testEnv struct {
	tripsService   *trips.Service
	invitesService *invites.Service
	store          *store.Store
}

// setupEnv initialises a real test database.
// Set TEST_DATABASE_URL to run these tests.
func setupEnv(t *testing.T) *testEnv {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set — skipping integration tests")
	}
	database, _, err := db.Connect(dsn)
	if err != nil {
		t.Fatalf("db.Connect: %v", err)
	}
	s := store.New(database)
	return &testEnv{
		tripsService:   trips.NewService(s),
		invitesService: invites.NewService(s, nil, nil), // nil bot — notifications disabled
		store:          s,
	}
}

// seedUser inserts a minimal user row for testing.
func seedUser(t *testing.T, db *gorm.DB) *models.User {
	t.Helper()
	u := &models.User{
		TelegramID: int64(uuid.New().ID()),
		Username:   uuid.New().String()[:10],
		FirstName:  "Test",
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seedUser: %v", err)
	}
	return u
}

// ── InviteUser ──────────────────────────────────────────────────────────────

func TestInviteUser_Success(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)

	trip, err := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{Title: "Test trip"})
	if err != nil {
		t.Fatalf("create trip: %v", err)
	}

	invite, err := env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")
	if err != nil {
		t.Fatalf("invite user: %v", err)
	}
	if invite.Status != "pending" {
		t.Errorf("status = %q, want pending", invite.Status)
	}
	if invite.InvitedUserID != invited.ID {
		t.Error("invited_user_id mismatch")
	}
}

func TestInviteUser_Duplicate(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)

	trip, _ := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{Title: "Dup test"})
	_, _ = env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")

	existing, err := env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")
	if err == nil || err != invites.ErrDuplicate {
		t.Fatalf("expected ErrDuplicate on second invite, got err=%v invite=%v", err, existing)
	}
	if existing == nil {
		t.Error("expected existing invite returned alongside ErrDuplicate")
	}
}

func TestInviteUser_NotOwner_PrivateTrip(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	other := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)

	visibility := "private"
	trip, _ := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{
		Title:      "Private",
		Visibility: visibility,
	})

	_, err := env.invitesService.InviteUser(ctx, other.ID, trip.ID, invited.Username, "")
	if err == nil {
		t.Fatal("expected ErrForbidden: non-member inviting to private trip")
	}
}

// ── RespondInvite ───────────────────────────────────────────────────────────

func TestRespondInvite_Accept_AddsMember(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)

	trip, _ := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{Title: "Accept test"})
	invite, _ := env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")

	if err := env.invitesService.RespondInvite(ctx, invited.ID, invite.ID, "accept", nil); err != nil {
		t.Fatalf("respond: %v", err)
	}

	isMember, err := env.store.IsTripMember(ctx, trip.ID, invited.ID)
	if err != nil {
		t.Fatalf("IsTripMember: %v", err)
	}
	if !isMember {
		t.Error("expected invited user to be a trip member after accepting")
	}
}

func TestRespondInvite_Decline_DoesNotAddMember(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)

	trip, _ := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{Title: "Decline test"})
	invite, _ := env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")

	if err := env.invitesService.RespondInvite(ctx, invited.ID, invite.ID, "decline", nil); err != nil {
		t.Fatalf("respond: %v", err)
	}

	isMember, _ := env.store.IsTripMember(ctx, trip.ID, invited.ID)
	if isMember {
		t.Error("declined user should NOT be a trip member")
	}
}

func TestRespondInvite_WrongUser(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)
	other := seedUser(t, env.store.DB)

	trip, _ := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{Title: "Wrong user"})
	invite, _ := env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")

	err := env.invitesService.RespondInvite(ctx, other.ID, invite.ID, "accept", nil)
	if err == nil {
		t.Fatal("expected ErrForbidden: wrong user responding to invite")
	}
}

// ── CancelInvite ────────────────────────────────────────────────────────────

func TestCancelInvite_OwnerCanCancel(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)

	trip, _ := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{Title: "Cancel test"})
	invite, _ := env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")

	if err := env.invitesService.CancelInvite(ctx, owner.ID, invite.ID); err != nil {
		t.Fatalf("cancel: %v", err)
	}
}

func TestCancelInvite_NonOwnerForbidden(t *testing.T) {
	env := setupEnv(t)
	ctx := context.Background()

	owner := seedUser(t, env.store.DB)
	invited := seedUser(t, env.store.DB)
	other := seedUser(t, env.store.DB)

	trip, _ := env.tripsService.CreateTrip(ctx, owner.ID, trips.CreateTripInput{Title: "Cancel perm"})
	invite, _ := env.invitesService.InviteUser(ctx, owner.ID, trip.ID, invited.Username, "")

	err := env.invitesService.CancelInvite(ctx, other.ID, invite.ID)
	if err == nil {
		t.Fatal("expected ErrForbidden: only trip owner can cancel")
	}
}
