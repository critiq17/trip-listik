package httpapi

import (
	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/httpapi/handlers"
	"github.com/critiq17/tripListik/internal/httpapi/middleware"
	"github.com/critiq17/tripListik/internal/invites"
	"github.com/critiq17/tripListik/internal/realtime"
	"github.com/critiq17/tripListik/internal/store"
	"github.com/critiq17/tripListik/internal/supabase"
	"github.com/critiq17/tripListik/internal/telegram"
	"github.com/critiq17/tripListik/internal/trips"
	"github.com/gofiber/fiber/v2"
)

func Register(v1 fiber.Router, cfg *config.Config, store *store.Store) {
	hub := realtime.NewHub()
	bot := telegram.NewBot(cfg.BotToken)

	// ── Services ──────────────────────────────────────────────────────────────
	tripsService := trips.NewService(store)
	invitesService := invites.NewService(store, bot, cfg)

	// ── Handlers ──────────────────────────────────────────────────────────────
	authHandler := &handlers.AuthHandler{Store: store, Cfg: cfg, Bot: bot}
	storageClient := supabase.NewStorageClient(cfg.SupabaseURL, cfg.SupabaseServiceKey)
	feedHandler := &handlers.FeedHandler{Store: store, Storage: storageClient, Bucket: cfg.SupabaseStorageBucket}
	tripsHandler := &handlers.TripsHandler{Service: tripsService, Storage: storageClient, Bucket: cfg.SupabaseStorageBucket}
	membersHandler := &handlers.MembersHandler{Store: store, Bot: bot}
	voteItemsHandler := &handlers.VoteItemsHandler{Store: store, Hub: hub}
	photosHandler := &handlers.PhotosHandler{Store: store, Storage: storageClient, Bucket: cfg.SupabaseStorageBucket, Hub: hub}
	streamHandler := &handlers.StreamHandler{Hub: hub}
	profileHandler := &handlers.ProfileHandler{Store: store}
	inboxHandler := &handlers.InboxHandler{Store: store}
	invitesHandler := &handlers.InvitesHandler{Service: invitesService, Store: store, Storage: storageClient, Bucket: cfg.SupabaseStorageBucket}
	usersHandler := &handlers.UsersHandler{Store: store}
	wishlistHandler := &handlers.WishlistHandler{Store: store}
	publicProfileHandler := &handlers.PublicProfileHandler{Store: store}
	geocodeHandler := &handlers.GeocodeHandler{}
	botHandler := &handlers.BotHandler{Store: store, Bot: bot, Cfg: cfg}

	// ── Telegram Bot Webhook (public, no auth, Telegram IPs only) ─────────────
	v1.Post("/webhook/telegram", botHandler.HandleUpdate)

	// ── Public ────────────────────────────────────────────────────────────────
	// OptionalAuth: these endpoints work anonymously, but resolve the viewer
	// when a token is present (private trip access, viewer_is_member, etc.).
	optionalAuth := middleware.OptionalAuth(cfg)
	v1.Post("/auth/telegram", authHandler.TelegramAuth)
	v1.Post("/auth/refresh", authHandler.Refresh)
	v1.Get("/feed", optionalAuth, feedHandler.Feed)
	v1.Get("/explore", optionalAuth, feedHandler.Explore)
	v1.Get("/trips/:id", optionalAuth, tripsHandler.GetTrip)
	v1.Get("/users/:id/profile", optionalAuth, publicProfileHandler.GetPublicProfile)
	v1.Get("/trips/:id/stream", middleware.RequireAuthQuery(cfg), streamHandler.TripStream)

	// ── Protected ─────────────────────────────────────────────────────────────
	protected := v1.Group("", middleware.RequireAuth(cfg))
	protected.Post("/trips", tripsHandler.CreateTrip)
	protected.Patch("/trips/:id", tripsHandler.UpdateTrip)
	protected.Delete("/trips/:id", tripsHandler.DeleteTrip)
	protected.Get("/trips", tripsHandler.ListMyTrips)

	// Members
	protected.Get("/trips/:id/members", membersHandler.ListMembers)
	protected.Post("/trips/:id/join", membersHandler.JoinTrip)
	protected.Get("/trips/:id/join/requests", membersHandler.ListJoinRequests)
	protected.Post("/trips/:id/join/approve", membersHandler.ApproveJoin)
	protected.Post("/trips/:id/join/reject", membersHandler.RejectJoin)
	protected.Delete("/trips/:id/members/:userId", membersHandler.RemoveMember)

	// Vote items (per spec: item-level voting on hotels, airlines, activities)
	protected.Get("/trips/:id/vote-items", voteItemsHandler.List)
	protected.Post("/trips/:id/vote-items", voteItemsHandler.Create)
	protected.Post("/trips/:id/vote-items/:itemId/vote", voteItemsHandler.Vote)
	protected.Delete("/trips/:id/vote-items/:itemId/vote", voteItemsHandler.Unvote)

	// Photos
	protected.Get("/trips/:id/photos", photosHandler.ListPhotos)
	protected.Post("/trips/:id/photos/presign", photosHandler.PresignUpload)
	protected.Post("/trips/:id/photos", photosHandler.CreatePhoto)
	protected.Delete("/trips/:id/photos/:photoId", photosHandler.DeletePhoto)

	// Invites
	protected.Post("/trips/:id/invite", invitesHandler.InviteUser)
	protected.Get("/trips/:id/invites", invitesHandler.ListTripInvites)
	protected.Get("/invites/:id", invitesHandler.GetInvite)
	protected.Post("/invites/:id/respond", invitesHandler.RespondInvite)
	protected.Delete("/invites/:id", invitesHandler.CancelInvite)

	// Search
	protected.Get("/users/search", usersHandler.Search)
	protected.Get("/geocode/search", geocodeHandler.Search)

	// Profile
	protected.Get("/me", profileHandler.Me)
	protected.Patch("/me", profileHandler.Update)
	protected.Get("/me/stats", profileHandler.Stats)
	protected.Get("/me/world", profileHandler.World)
	protected.Get("/me/map", profileHandler.Map)

	// Wishlist
	protected.Get("/me/wishlist", wishlistHandler.List)
	protected.Post("/me/wishlist", wishlistHandler.Create)
	protected.Delete("/me/wishlist/:id", wishlistHandler.Delete)

	// Inbox / Notifications
	protected.Get("/inbox", inboxHandler.ListInbox)
	protected.Get("/inbox/unread", inboxHandler.UnreadCount)
	protected.Post("/inbox/read-all", inboxHandler.MarkAllRead)

	// Referrals
	protected.Get("/me/referrals", profileHandler.Referrals)

	// Friends
	protected.Get("/me/friends", profileHandler.Friends)
}
