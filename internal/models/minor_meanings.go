package models

import (
	"database/sql"

	"github.com/ilbagatto/tarot-api/internal/utils"
)

// MeaningMinor represents an interpretation of a minor arcana card
type MeaningMinor struct {
	ID       int64           `json:"id"`
	Suit     int64           `json:"suit"`
	Rank     int64           `json:"rank"`
	Position MeaningPosition `json:"position"`
	Source   int64           `json:"source"`
	Meaning  string          `json:"meaning"`
}

// MeaningMinorInput is used for creating or updating MeaningMinor records
type MeaningMinorInput struct {
	Suit     int64           `json:"suit" example:"1"`
	Rank     int64           `json:"rank" example:"5"`
	Position MeaningPosition `json:"position" example:"straight"`
	Source   int64           `json:"source" example:"1"`
	Meaning  string          `json:"meaning" example:"Active communication and drive"`
}

// ListMinorMeaning returns all MeaningMinor entries for given suit, name, position and source
func ListMinorMeanings(db *sql.DB, filters map[string]any) ([]MeaningMinor, error) {
	query := `
	SELECT id, suit, rank, position, source, meaning
	FROM meaning_minor
	`

	whereClause, args := utils.BuildWhereClause(filters, 1)
	query += " " + whereClause + " ORDER BY position"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meanings []MeaningMinor
	for rows.Next() {
		var m MeaningMinor
		if err := rows.Scan(&m.ID, &m.Suit, &m.Rank, &m.Position, &m.Source, &m.Meaning); err != nil {
			return nil, err
		}
		meanings = append(meanings, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return meanings, nil
}

// GetMinorMeaningByID retrieves a MeaningMinor by its unique ID
func GetMinorMeaningByID(db *sql.DB, id int64) (*MeaningMinor, error) {
	const query = `
	SELECT id, suit, rank, position, source, meaning
	FROM meaning_minor
	WHERE id = $1`
	var m MeaningMinor
	if err := db.QueryRow(query, id).Scan(
		&m.ID, &m.Suit, &m.Rank, &m.Position, &m.Source, &m.Meaning,
	); err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateMinorMeaning inserts a new record into meaning_minor
func CreateMinorMeaning(db *sql.DB, input MeaningMinorInput) (*int64, error) {
	const query = `
	INSERT INTO meaning_minor (suit, rank, position, source, meaning)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	var id int64
	if err := db.QueryRow(query, input.Suit, input.Rank, input.Position, input.Source, input.Meaning).Scan(&id); err != nil {
		return nil, err
	}
	return &id, nil
}

// UpdateMinorMeaning updates an existing record by ID
func UpdateMinorMeaning(db *sql.DB, id int64, input MeaningMinorInput) (*MeaningMinor, error) {
	const query = `
	UPDATE meaning_minor
	SET suit = $1, rank = $2, position = $3, source = $4, meaning = $5
	WHERE id = $6`
	res, err := db.Exec(query, input.Suit, input.Rank, input.Position, input.Source, input.Meaning, id)
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

	updated := &MeaningMinor{
		ID: id, Suit: input.Suit,
		Rank:     input.Rank,
		Position: input.Position,
		Source:   input.Source,
		Meaning:  input.Meaning,
	}
	return updated, nil
}

// DeleteMinorMeaning removes a record from the meaning_minor table by ID
func DeleteMinorMeaning(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM meaning_minor WHERE id = $1`, id)
	return err
}
