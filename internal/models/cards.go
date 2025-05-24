package models

import "database/sql"

type Card struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name" example:"The Fool"`
	DeckID    int64        `json:"deck" example:"1"`
	Image     *string      `json:"image,omitempty"`     // Full URL
	Thumbnail *string      `json:"thumbnail,omitempty"` // Full URL
	Meanings  []MeaningRef `json:"meanings,omitempty"`
}

func updateCard(tx *sql.Tx, deckID int64, id int64) error {
	res, err := tx.Exec("UPDATE card SET deck = $1 WHERE id = $2", deckID, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func insertCard(tx *sql.Tx, deckID int64, arcana string) (*int64, error) {
	const insertCard = `
		INSERT INTO card (deck, arcana)
		VALUES ($1, $2)
		RETURNING id`
	var cardID int64
	if err := tx.QueryRow(insertCard, deckID, arcana).Scan(&cardID); err != nil {
		return nil, err
	}
	return &cardID, nil
}
