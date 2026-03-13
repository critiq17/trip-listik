package models

import "time"

// VoteItem represents a hotel, airline, activity or other option that trip members can vote on.
type VoteItem struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TripID      string    `gorm:"type:uuid;index" json:"trip_id"`
	Category    string    `gorm:"type:text" json:"category"` // HOTEL, AIRLINE, ACTIVITY, OTHER
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	AddedBy     string    `gorm:"type:uuid" json:"added_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// VoteItemVote represents a user's vote on a specific VoteItem.
type VoteItemVote struct {
	ItemID    string    `gorm:"type:uuid;primaryKey" json:"item_id"`
	UserID    string    `gorm:"type:uuid;primaryKey" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (VoteItem) TableName() string     { return "vote_items" }
func (VoteItemVote) TableName() string { return "vote_item_votes" }
