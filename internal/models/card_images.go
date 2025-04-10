package models

import "database/sql"

// CardImage represents an image associated with a tarot card
type CardImage struct {
	CardID int64  `json:"card_id"`
	Path   string `json:"path"`
}

func GetCardImageByCardID(db *sql.DB, cardID int64) (*CardImage, error) {
	const query = `SELECT card, path FROM card_image WHERE card = $1`

	var img CardImage
	err := db.QueryRow(query, cardID).Scan(&img.CardID, &img.Path)
	if err != nil {
		return nil, err
	}
	return &img, nil
}
