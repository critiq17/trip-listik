package handlers

import (
	"strings"

	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/critiq17/tripListik/internal/supabase"
)

func normalizeStoredPhotoURL(storage *supabase.StorageClient, bucket, raw string) string {
	cleanRaw := strings.TrimSpace(raw)
	if cleanRaw == "" {
		return ""
	}
	if storage == nil {
		return cleanRaw
	}
	return storage.CanonicalPublicURL(bucket, cleanRaw)
}

func normalizeTripCoverPhoto(storage *supabase.StorageClient, bucket string, trip *models.Trip) {
	if trip == nil {
		return
	}
	trip.CoverPhotoURL = normalizeStoredPhotoURL(storage, bucket, trip.CoverPhotoURL)
}

func normalizeTripCoverPhotos(storage *supabase.StorageClient, bucket string, trips []models.Trip) {
	for i := range trips {
		normalizeTripCoverPhoto(storage, bucket, &trips[i])
	}
}

func normalizeTripFeedItems(storage *supabase.StorageClient, bucket string, items []TripFeedItem) {
	for i := range items {
		items[i].CoverPhotoURL = normalizeStoredPhotoURL(storage, bucket, items[i].CoverPhotoURL)
	}
}

func normalizeInviteViews(storage *supabase.StorageClient, bucket string, items []store.InviteView) {
	for i := range items {
		items[i].TripCover = normalizeStoredPhotoURL(storage, bucket, items[i].TripCover)
	}
}

func normalizeTripPhoto(storage *supabase.StorageClient, bucket string, photo *models.TripPhoto) {
	if photo == nil {
		return
	}

	if storage != nil && photo.StoragePath == "" {
		photo.StoragePath = storage.NormalizeObjectPath(bucket, photo.URL)
	}

	source := strings.TrimSpace(photo.StoragePath)
	if source == "" {
		source = photo.URL
	}
	photo.URL = normalizeStoredPhotoURL(storage, bucket, source)
}

func normalizeTripPhotos(storage *supabase.StorageClient, bucket string, photos []models.TripPhoto) {
	for i := range photos {
		normalizeTripPhoto(storage, bucket, &photos[i])
	}
}
