package handlers

import (
	"strings"

	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/critiq17/tripListik/internal/supabase"
)

const photoURLExpirySeconds = 7 * 24 * 60 * 60

func normalizeStoredPhotoURL(storage *supabase.StorageClient, bucket, raw string) string {
	cleanRaw := strings.TrimSpace(raw)
	if cleanRaw == "" {
		return ""
	}
	if storage == nil {
		return cleanRaw
	}

	cleanPath := storage.NormalizeObjectPath(bucket, cleanRaw)
	if cleanPath == "" {
		return storage.CanonicalPublicURL(bucket, cleanRaw)
	}

	signedURL, err := storage.CreateSignedObjectURL(bucket, cleanPath, photoURLExpirySeconds)
	if err == nil && signedURL != "" {
		return signedURL
	}

	publicURL := storage.PublicObjectURL(bucket, cleanPath)
	if publicURL != "" {
		return publicURL
	}

	return cleanPath
}

func normalizeTripCoverPhoto(storage *supabase.StorageClient, bucket string, trip *models.Trip) {
	if trip == nil {
		return
	}
	trip.CoverPhotoURL = normalizeStoredPhotoURL(storage, bucket, trip.CoverPhotoURL)
}

func normalizeTripCoverPhotos(storage *supabase.StorageClient, bucket string, trips []models.Trip) {
	signedURLs := signedPhotoURLMap(storage, bucket, tripCoverPaths(trips))
	for i := range trips {
		normalizeTripCoverPhotoWithMap(storage, bucket, &trips[i], signedURLs)
	}
}

func normalizeTripFeedItems(storage *supabase.StorageClient, bucket string, items []TripFeedItem) {
	signedURLs := signedPhotoURLMap(storage, bucket, tripFeedPaths(items))
	for i := range items {
		items[i].CoverPhotoURL = normalizeStoredPhotoURLWithMap(storage, bucket, items[i].CoverPhotoURL, signedURLs)
	}
}

func normalizeInviteViews(storage *supabase.StorageClient, bucket string, items []store.InviteView) {
	signedURLs := signedPhotoURLMap(storage, bucket, inviteCoverPaths(items))
	for i := range items {
		items[i].TripCover = normalizeStoredPhotoURLWithMap(storage, bucket, items[i].TripCover, signedURLs)
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
	signedURLs := signedPhotoURLMap(storage, bucket, tripPhotoPaths(photos))
	for i := range photos {
		normalizeTripPhotoWithMap(storage, bucket, &photos[i], signedURLs)
	}
}

func normalizeTripCoverPhotoWithMap(storage *supabase.StorageClient, bucket string, trip *models.Trip, signedURLs map[string]string) {
	if trip == nil {
		return
	}
	trip.CoverPhotoURL = normalizeStoredPhotoURLWithMap(storage, bucket, trip.CoverPhotoURL, signedURLs)
}

func normalizeTripPhotoWithMap(storage *supabase.StorageClient, bucket string, photo *models.TripPhoto, signedURLs map[string]string) {
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
	photo.URL = normalizeStoredPhotoURLWithMap(storage, bucket, source, signedURLs)
}

func normalizeStoredPhotoURLWithMap(storage *supabase.StorageClient, bucket, raw string, signedURLs map[string]string) string {
	cleanRaw := strings.TrimSpace(raw)
	if cleanRaw == "" || storage == nil {
		return normalizeStoredPhotoURL(storage, bucket, raw)
	}

	cleanPath := storage.NormalizeObjectPath(bucket, cleanRaw)
	if cleanPath != "" {
		if signedURL := strings.TrimSpace(signedURLs[cleanPath]); signedURL != "" {
			return signedURL
		}
	}

	return normalizeStoredPhotoURL(storage, bucket, cleanRaw)
}

func signedPhotoURLMap(storage *supabase.StorageClient, bucket string, paths []string) map[string]string {
	if storage == nil || len(paths) == 0 {
		return nil
	}

	signedURLs, err := storage.CreateSignedObjectURLs(bucket, paths, photoURLExpirySeconds)
	if err != nil {
		return nil
	}

	return signedURLs
}

func tripCoverPaths(trips []models.Trip) []string {
	paths := make([]string, 0, len(trips))
	for _, trip := range trips {
		if trip.CoverPhotoURL != "" {
			paths = append(paths, trip.CoverPhotoURL)
		}
	}
	return paths
}

func tripFeedPaths(items []TripFeedItem) []string {
	paths := make([]string, 0, len(items))
	for _, item := range items {
		if item.CoverPhotoURL != "" {
			paths = append(paths, item.CoverPhotoURL)
		}
	}
	return paths
}

func inviteCoverPaths(items []store.InviteView) []string {
	paths := make([]string, 0, len(items))
	for _, item := range items {
		if item.TripCover != "" {
			paths = append(paths, item.TripCover)
		}
	}
	return paths
}

func tripPhotoPaths(photos []models.TripPhoto) []string {
	paths := make([]string, 0, len(photos))
	for _, photo := range photos {
		if photo.StoragePath != "" {
			paths = append(paths, photo.StoragePath)
			continue
		}
		if photo.URL != "" {
			paths = append(paths, photo.URL)
		}
	}
	return paths
}
