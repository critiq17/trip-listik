package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TelegramID int64     `gorm:"uniqueIndex" json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	PhotoURL   string    `json:"photo_url"`
	Bio        string    `json:"bio"`
	IsPublic   bool      `json:"is_public"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Trip struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OwnerID       uuid.UUID  `gorm:"type:uuid;index" json:"owner_id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	StartDate     *time.Time `gorm:"type:date" json:"start_date"`
	EndDate       *time.Time `gorm:"type:date" json:"end_date"`
	Visibility    string     `json:"visibility"`
	Status        string     `json:"status"`
	CountryCode   string     `json:"country_code"`
	City          string     `json:"city"`
	CoverPhotoURL string     `json:"cover_photo_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"-"`
}

type TripMember struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID   uuid.UUID `gorm:"type:uuid;index" json:"trip_id"`
	UserID   uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type TripJoinRequest struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID    uuid.UUID `gorm:"type:uuid;index" json:"trip_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TripInvite struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID          uuid.UUID `gorm:"type:uuid;index" json:"trip_id"`
	InvitedUserID   uuid.UUID `gorm:"type:uuid;index" json:"invited_user_id"`
	InvitedByUserID uuid.UUID `gorm:"type:uuid;index" json:"invited_by_user_id"`
	Status          string    `json:"status"`
	Comment         *string   `json:"comment"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type TripVote struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID    uuid.UUID `gorm:"type:uuid;index" json:"trip_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Vote      int16     `json:"vote"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TripComment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID    uuid.UUID `gorm:"type:uuid;index" json:"trip_id"`
	UserID    uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TripPhoto struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID      uuid.UUID `gorm:"type:uuid;index" json:"trip_id"`
	UserID      uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	StoragePath string    `json:"storage_path"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserCountry struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	CountryCode  string    `json:"country_code"`
	FirstVisitAt time.Time `json:"first_visit_at"`
}

type UserStats struct {
	UserID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	TotalTrips       int       `json:"total_trips"`
	CountriesVisited int       `json:"countries_visited"`
	CitiesVisited    int       `json:"cities_visited"`
	TripsWithFriends int       `json:"trips_with_friends"`
	SoloTrips        int       `json:"solo_trips"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UserWishlistItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	CountryCode string    `json:"country_code"`
	City        string    `json:"city"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
}

func (UserWishlistItem) TableName() string {
	return "user_wishlist"
}

type Notification struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;index" json:"user_id"`
	Type      string     `json:"type"`
	TripID    *uuid.UUID `gorm:"type:uuid" json:"trip_id"`
	ActorID   *uuid.UUID `gorm:"type:uuid" json:"actor_id"`
	Payload   []byte     `gorm:"type:jsonb" json:"payload"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type RefreshToken struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;index" json:"user_id"`
	TokenHash  string     `gorm:"uniqueIndex" json:"-"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}
