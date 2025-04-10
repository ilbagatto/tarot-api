package models

import (
	"database/sql"

	"github.com/ilbagatto/tarot-api/internal/utils"
)

// CardMajor represents a Major Arcana card
type CardMajor struct {
	Card
	Number  int    `json:"number" example:"0"`
	Name    string `json:"name" example:"The Fool"`
	OrgName string `json:"orgname,omitempty" example:"Le Mat"`
}

// CardMajorInput is used for creating or updating Major Arcana cards
type CardMajorInput struct {
	DeckID  int64  `json:"deck" example:"1"`
	Number  int    `json:"number" example:"0"`
	Name    string `json:"name" example:"The Fool"`
	OrgName string `json:"orgname,omitempty" example:"Le Mat"`
}

// ListMajorCards retrieves all Major Arcana cards for a given deck
func ListMajorCards(db *sql.DB, deckID int64) ([]CardMajor, error) {
	const query = `
		SELECT c.id, c.deck, m.number, m.name, m.orgname
		FROM card c
		JOIN card_major m ON m.card = c.id
		WHERE c.deck = $1
		ORDER BY m.number`

	rows, err := db.Query(query, deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []CardMajor
	for rows.Next() {
		var card CardMajor
		if err := rows.Scan(&card.ID, &card.DeckID, &card.Number, &card.Name, &card.OrgName); err != nil {
			return nil, err
		}

		img, err := GetCardImageByCardID(db, card.ID)
		if err == nil {
			card.Image = utils.GetCardImage(img.Path)
		}

		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cards, nil
}

// GetMajorCardByID retrieves a Major Arcana card by its ID
func GetMajorCardByID(db *sql.DB, id int64) (*CardMajor, error) {
	var query = `
		SELECT c.id, c.deck, m.number, m.name, m.orgname
		FROM card c
		JOIN card_major m ON m.card = c.id
		WHERE c.id = $1`

	var card CardMajor
	if err := db.QueryRow(query, id).Scan(
		&card.ID, &card.DeckID, &card.Number, &card.Name, &card.OrgName,
	); err != nil {
		return nil, err
	}

	img, err := GetCardImageByCardID(db, card.ID)
	if err == nil {
		card.Image = utils.GetCardImage(img.Path)
	}

	// Load related meanings
	query = `
		SELECT m.id, m.position, m.source
		FROM meaning_major m
		WHERE m.number = $1
		AND m.source IN (
			SELECT source FROM deck_source WHERE deck = $2
		)
		ORDER BY m.source, m.number, m.position
	`
	rows, err := db.Query(query, card.Number, card.DeckID)
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

// CreateMajorCard inserts a new Major Arcana card into the database
func CreateMajorCard(db *sql.DB, input CardMajorInput) (*int64, error) {
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

	cardID, err := insertCard(tx, input.DeckID, "major")
	if err != nil {
		return nil, err
	}

	const insertMajor = `
		INSERT INTO card_major (card, number, name, orgname)
		VALUES ($1, $2, $3, $4)`
	if _, err := tx.Exec(insertMajor, &cardID, input.Number, input.Name, input.OrgName); err != nil {
		return nil, err
	}

	return cardID, nil
}

// UpdateMajorCard updates an existing Major Arcana card
func UpdateMajorCard(db *sql.DB, id int64, input CardMajorInput) (*CardMajor, error) {
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

	if err := updateCard(tx, input.DeckID, id); err != nil {
		return nil, err
	}

	const updateMajor = `
		UPDATE card_major
		SET number = $1, name = $2, orgname = $3
		WHERE card = $4`
	res, err := tx.Exec(updateMajor, input.Number, input.Name, input.OrgName, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	updated := &CardMajor{
		Card: Card{
			ID:     id,
			Name:   input.Name,
			DeckID: input.DeckID,
			Image:  nil,
		},
		Number:  input.Number,
		Name:    input.Name,
		OrgName: input.OrgName,
	}

	img, err := GetCardImageByCardID(db, id)
	if err == nil {
		updated.Image = utils.GetCardImage(img.Path)
	}

	return updated, nil
}

// DeleteMajorCard deletes a Major Arcana card from the database
func DeleteMajorCard(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM card WHERE id = $1", id)
	return err
}
