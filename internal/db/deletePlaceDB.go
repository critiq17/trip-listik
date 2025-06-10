package db

import "context"

func DeletePlaceFromDB(chatID int64, place string) error {
	_, err := DB.Exec(context.Background(),
		"DELETE FROM user-wishlist WHERE chat_id = $1 AND place = $2", chatID, place)
	return err
}
