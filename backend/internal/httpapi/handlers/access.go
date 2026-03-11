package handlers

import (
	"context"
	"time"

	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/google/uuid"
)

func ensureTripAccess(ctx context.Context, store *store.Store, tripID uuid.UUID, userID *uuid.UUID) (*models.Trip, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	trip, err := store.GetTripByID(ctx, tripID)
	if err != nil || trip == nil {
		return trip, false, err
	}

	if trip.Visibility == "public" {
		return trip, true, nil
	}
	if userID == nil {
		return trip, false, nil
	}

	isMember, err := store.IsTripMember(ctx, tripID, *userID)
	if err != nil {
		return nil, false, err
	}
	if !isMember {
		return trip, false, nil
	}

	return trip, true, nil
}
