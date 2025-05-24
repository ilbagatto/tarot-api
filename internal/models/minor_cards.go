package models

import (
	"database/sql"

	"github.com/ilbagatto/tarot-api/internal/utils"
)

// CardMinor represents a Minor Arcana card
type CardMinor struct {
	Card
	SuitID int64 `json:"suit" example:"1"`
	RankID int64 `json:"rank" example:"10"`
}

// CardMinorInput is used to create or update Minor Arcana cards
type CardMinorInput struct {
	DeckID int64 `json:"deck" example:"1"`
	SuitID int64 `json:"suit" example:"1"`
	RankID int64 `json:"rank" example:"10"`
}

// ListMinorCards retrieves all Minor Arcana cards for a given deck
func ListMinorCards(db *sql.DB, deckID int64) ([]CardMinor, error) {
	const query = `
	SELECT c.id, CONCAT(r.name, ' ', s.genitive) AS name, c.deck, m.suit, m.rank
	FROM card_minor m
	JOIN card c ON c.id = m.card
	JOIN rank r ON r.id = m.rank
	JOIN suit s ON s.id = m.suit
	WHERE c.deck = $1
	ORDER BY m.suit, m.rank`

	rows, err := db.Query(query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []CardMinor
	for rows.Next() {
		var card CardMinor
		if err := rows.Scan(&card.ID, &card.Name, &card.DeckID, &card.SuitID, &card.RankID); err != nil {
			return nil, err
		}

		img, err := GetCardImageByCardID(db, card.ID)
		if err == nil {
			card.Image = utils.GetImageURL(img.Path, false)
			card.Thumbnail = utils.GetImageURL(img.Path, true)
		}

		cards = append(cards, card)
	}
	return cards, rows.Err()
}

// GetMinorCardByID retrieves a Minor Arcana card by its ID
func GetMinorCardByID(db *sql.DB, id int64) (*CardMinor, error) {
	var query = `
	SELECT c.id, CONCAT(r.name, ' ', s.genitive) AS name, c.deck, m.suit, m.rank
	FROM card_minor m
	JOIN card c ON c.id = m.card
	JOIN rank r ON r.id = m.rank
	JOIN suit s ON s.id = m.suit	
	WHERE c.id = $1`

	var card CardMinor
	if err := db.QueryRow(query, id).Scan(
		&card.ID, &card.Name, &card.DeckID, &card.SuitID, &card.RankID,
	); err != nil {
		return nil, err
	}

	img, err := GetCardImageByCardID(db, card.ID)
	if err == nil {
		card.Image = utils.GetImageURL(img.Path, false)
		card.Thumbnail = utils.GetImageURL(img.Path, true)
	}

	// Load related meanings
	query = `
		SELECT m.id, m.position, m.source
		FROM meaning_minor m
		WHERE m.suit = $1 AND m.rank = $2
		AND m.source IN (
			SELECT source FROM deck_source WHERE deck = $3
		)
		ORDER BY m.source, m.suit, m.rank, m.position
	`
	rows, err := db.Query(query, card.SuitID, card.RankID, card.DeckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m MeaningRef
		if err := rows.Scan(&m.ID, &m.Position, &m.SourceID); err != nil {
			return nil, err
		}
		card.Meanings = append(card.Meanings, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &card, nil
}

// CreateMinorCard inserts a new Minor Arcana card into the database
func CreateMinorCard(db *sql.DB, input CardMinorInput) (*int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Insert into card
	cardID, err := insertCard(tx, input.DeckID, "minor")
	if err != nil {
		return nil, err
	}
	// Insert into card_minor
	const insertMinor = `
		INSERT INTO card_minor (card, suit, rank)
		VALUES ($1, $2, $3)`
	if _, err = tx.Exec(insertMinor, &cardID, input.SuitID, input.RankID); err != nil {
		return nil, err
	}

	return cardID, nil
}

// UpdateMinorCard updates an existing Minor Arcana card
func UpdateMinorCard(db *sql.DB, id int64, input CardMinorInput) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err := updateCard(tx, input.DeckID, id); err != nil {
		return err
	}

	const updateMinor = `
		UPDATE card_minor SET suit = $1, rank = $2 WHERE card = $3`
	res, err := tx.Exec(updateMinor, input.SuitID, input.RankID, id)
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

func DeleteMinorCard(db *sql.DB, id int64) error {
	const query = `DELETE FROM card WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
