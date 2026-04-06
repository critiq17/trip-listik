package handlers

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/httpapi/validate"
	"github.com/critiq17/tripListik/internal/realtime"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/store/models"
	"github.com/critiq17/tripListik/internal/supabase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PhotosHandler struct {
	Store   *store.Store
	Storage *supabase.StorageClient
	Bucket  string
	Hub     *realtime.Hub
}

type createPhotoRequest struct {
	StoragePath string `json:"storage_path" validate:"required"`
	URL         string `json:"url" validate:"omitempty,url"`
}

type presignRequest struct {
	FileName    string `json:"file_name" validate:"required"`
	ContentType string `json:"content_type"`
}

func (h *PhotosHandler) ListPhotos(c *fiber.Ctx) error {
	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	_, allowed, err := ensureTripAccess(c.Context(), h.Store, tripID, &userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}
	limit := c.QueryInt("limit", 50)
	if limit <= 0 || limit > 100 {
		limit = 50
	}
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

	photos, err := h.Store.ListPhotos(ctx, tripID, limit, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load photos")
	}
	normalizeTripPhotos(h.Storage, h.Bucket, photos)

	return c.JSON(fiber.Map{"items": photos})
}

func (h *PhotosHandler) CreatePhoto(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	_, allowed, err := ensureTripAccess(c.Context(), h.Store, tripID, &userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	var req createPhotoRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "storage_path is required")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	storagePath := strings.TrimSpace(req.StoragePath)
	if h.Storage != nil {
		storagePath = h.Storage.NormalizeObjectPath(h.Bucket, storagePath)
	}
	if storagePath == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid storage_path")
	}

	photoURL := strings.TrimSpace(req.URL)
	if h.Storage != nil {
		canonicalURL := h.Storage.PublicObjectURL(h.Bucket, storagePath)
		if canonicalURL != "" {
			photoURL = canonicalURL
		}
	}
	if photoURL == "" {
		return fiber.NewError(fiber.StatusBadRequest, "url is required")
	}

	photo := &models.TripPhoto{
		TripID:      tripID,
		UserID:      userID,
		StoragePath: storagePath,
		URL:         photoURL,
	}

	if err := h.Store.CreatePhoto(ctx, photo); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create photo")
	}
	normalizeTripPhoto(h.Storage, h.Bucket, photo)

	if h.Hub != nil {
		h.Hub.Publish("trip:"+tripID.String(), realtime.Event{
			Type: "photo_created",
			Data: photo,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(photo)
}

func (h *PhotosHandler) DeletePhoto(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	photoID, err := uuid.Parse(c.Params("photoId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid photo id")
	}

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	if err := h.Store.DeletePhoto(ctx, photoID, userID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete photo")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PhotosHandler) PresignUpload(c *fiber.Ctx) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}

	tripID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid trip id")
	}

	_, allowed, err := ensureTripAccess(c.Context(), h.Store, tripID, &userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load trip")
	}
	if !allowed {
		return fiber.NewError(fiber.StatusForbidden, "forbidden")
	}

	var req presignRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validate.Struct(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file_name is required")
	}

	ext := filepath.Ext(req.FileName)
	objectID := uuid.New().String() + ext
	objectPath := "trip_photos/" + tripID.String() + "/" + objectID

	upload, err := h.Storage.CreateSignedUploadURL(h.Bucket, objectPath, 3600)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create signed upload url")
	}

	// Supabase returns a relative signedUrl (e.g. /storage/v1/object/upload/sign/...)
	// Prepend the base URL so the client can upload directly to Supabase.
	signedURL := upload.SignedURL
	if len(signedURL) > 0 && signedURL[0] == '/' {
		signedURL = h.Storage.BaseURL + signedURL
	}

	// Supabase may echo back the path with the bucket name prepended
	// (e.g. "trip-photos/trip_photos/uuid/file.jpg"). Strip the bucket prefix
	// so the frontend receives just the object path ("trip_photos/uuid/file.jpg").
	storagePath := upload.Path
	if storagePath == "" {
		storagePath = objectPath
	}
	if h.Storage != nil {
		storagePath = h.Storage.NormalizeObjectPath(h.Bucket, storagePath)
	}

	return c.JSON(fiber.Map{
		"signed_url": signedURL,
		"token":      upload.Token,
		"path":       storagePath,
		"public_url": h.Storage.PublicObjectURL(h.Bucket, storagePath),
	})
}
