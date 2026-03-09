package httpapi

import (
	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/httpapi/handlers"
	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/realtime"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/supabase"
	"github.com/gofiber/fiber/v2"
)

func Register(v1 fiber.Router, cfg *config.Config, store *store.Store) {
	hub := realtime.NewHub()
	authHandler := &handlers.AuthHandler{Store: store, Cfg: cfg}
	feedHandler := &handlers.FeedHandler{Store: store}
	tripsHandler := &handlers.TripsHandler{Store: store}
	membersHandler := &handlers.MembersHandler{Store: store}
	votesHandler := &handlers.VotesHandler{Store: store, Hub: hub}
	commentsHandler := &handlers.CommentsHandler{Store: store, Hub: hub}
	storageClient := supabase.NewStorageClient(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	photosHandler := &handlers.PhotosHandler{Store: store, Storage: storageClient, Bucket: cfg.SupabaseStorageBucket, Hub: hub}
	streamHandler := &handlers.StreamHandler{Hub: hub}
	profileHandler := &handlers.ProfileHandler{Store: store}
	inboxHandler := &handlers.InboxHandler{Store: store}

	v1.Post("/auth/telegram", authHandler.TelegramAuth)
	v1.Get("/feed", feedHandler.Feed)
	v1.Get("/explore", feedHandler.Explore)
	v1.Get("/trips/:id", tripsHandler.GetTrip)
	v1.Get("/trips/:id/stream", middleware.RequireAuthQuery(cfg), streamHandler.TripStream)

	protected := v1.Group("", middleware.RequireAuth(cfg))
	protected.Post("/trips", tripsHandler.CreateTrip)
	protected.Patch("/trips/:id", tripsHandler.UpdateTrip)
	protected.Get("/trips", tripsHandler.ListMyTrips)

	protected.Get("/trips/:id/members", membersHandler.ListMembers)
	protected.Post("/trips/:id/join", membersHandler.JoinTrip)
	protected.Get("/trips/:id/join/requests", membersHandler.ListJoinRequests)
	protected.Post("/trips/:id/join/approve", membersHandler.ApproveJoin)
	protected.Post("/trips/:id/join/reject", membersHandler.RejectJoin)
	protected.Delete("/trips/:id/members/:userId", membersHandler.RemoveMember)

	protected.Post("/trips/:id/votes", votesHandler.CastVote)
	protected.Get("/trips/:id/votes", votesHandler.GetVotes)

	protected.Get("/trips/:id/comments", commentsHandler.ListComments)
	protected.Post("/trips/:id/comments", commentsHandler.CreateComment)

	protected.Get("/trips/:id/photos", photosHandler.ListPhotos)
	protected.Post("/trips/:id/photos/presign", photosHandler.PresignUpload)
	protected.Post("/trips/:id/photos", photosHandler.CreatePhoto)
	protected.Delete("/trips/:id/photos/:photoId", photosHandler.DeletePhoto)

	protected.Get("/me", profileHandler.Me)
	protected.Get("/me/stats", profileHandler.Stats)
	protected.Get("/me/world", profileHandler.World)

	protected.Get("/inbox", inboxHandler.ListInbox)
}
