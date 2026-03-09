package handlers

import (
	"context"
	"path/filepath"
	"time"

	"github.com/critiq17/tripListik/internal/httpapi/middleware"
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
	StoragePath string `json:"storage_path"`
	URL         string `json:"url"`
}

type presignRequest struct {
	FileName    string `json:"file_name"`
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

	_, allowed, err := ensureTripAccess(h.Store, tripID, &userID)
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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	photos, err := h.Store.ListPhotos(ctx, tripID, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to load photos")
	}

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

	_, allowed, err := ensureTripAccess(h.Store, tripID, &userID)
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
	if req.StoragePath == "" || req.URL == "" {
		return fiber.NewError(fiber.StatusBadRequest, "storage_path and url are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	photo := &models.TripPhoto{
		TripID:      tripID,
		UserID:      userID,
		StoragePath: req.StoragePath,
		URL:         req.URL,
	}

	if err := h.Store.CreatePhoto(ctx, photo); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create photo")
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

	_, allowed, err := ensureTripAccess(h.Store, tripID, &userID)
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
	if req.FileName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "file_name is required")
	}

	ext := filepath.Ext(req.FileName)
	objectID := uuid.New().String() + ext
	objectPath := "trip_photos/" + tripID.String() + "/" + objectID

	upload, err := h.Storage.CreateSignedUploadURL(h.Bucket, objectPath, 3600)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create signed upload url")
	}

	return c.JSON(fiber.Map{
		"signed_url": upload.SignedURL,
		"token":      upload.Token,
		"path":       upload.Path,
	})
}
