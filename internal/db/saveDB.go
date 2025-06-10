package db

import (
	"context"
)

func SavePlaceToDB(chatID int64, place string) error {
	_, err := DB.Exec(context.Background(),
		"INSERT INTO wishlist (chat_id, place) VALUES ($1, $2)", chatID, place)
	return err
}
