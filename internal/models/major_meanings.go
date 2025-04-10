package models

import (
	"database/sql"

	"github.com/ilbagatto/tarot-api/internal/utils"
)

// MeaningMajor represents an interpretation of a major arcana card
type MeaningMajor struct {
	ID       int64           `json:"id"`
	Number   int             `json:"number"`   // major arcana card number (0–21)
	Position MeaningPosition `json:"position"` // straight or reversed
	Source   int64           `json:"source"`
	Meaning  string          `json:"meaning"`
}

// MeaningMajorInput is used to create or update a MeaningMajor
type MeaningMajorInput struct {
	Number   int             `json:"number" example:"5"`
	Position MeaningPosition `json:"position" example:"straight"`
	Source   int64           `json:"source" example:"1"`
	Meaning  string          `json:"meaning" example:"Spiritual wisdom and intuition"`
}

// ListMajorMeaning returns all MeaningMajor entries for given number and source
func ListMajorMeanings(db *sql.DB, filters map[string]any) ([]MeaningMajor, error) {
	query := `
		SELECT id, number, position, source, meaning
		FROM meaning_major
	`

	whereClause, args := utils.BuildWhereClause(filters, 1)
	query += " " + whereClause + " ORDER BY position"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meanings []MeaningMajor
	for rows.Next() {
		var m MeaningMajor
		if err := rows.Scan(&m.ID, &m.Number, &m.Position, &m.Source, &m.Meaning); err != nil {
			return nil, err
		}
		meanings = append(meanings, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return meanings, nil
}

// GetMajorMeaningByID retrieves a MeaningMajor by its unique ID
func GetMajorMeaningByID(db *sql.DB, id int64) (*MeaningMajor, error) {
	const query = `
	SELECT id, number, position, source, meaning
	FROM meaning_major
	WHERE id = $1`
	var m MeaningMajor
	if err := db.QueryRow(query, id).Scan(&m.ID, &m.Number, &m.Position, &m.Source, &m.Meaning); err != nil {
		return nil, err
	}
	return &m, nil
}

// CreateMeaningMajor inserts a new MeaningMajor record
func CreateMeaningMajor(db *sql.DB, input MeaningMajorInput) (*int64, error) {
	const query = `
	INSERT INTO meaning_major (number, position, source, meaning)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var id int64
	if err := db.QueryRow(query, input.Number, input.Position, input.Source, input.Meaning).Scan(&id); err != nil {
		return nil, err
	}
	return &id, nil
}

// UpdateMajorMeaning updates an existing record by ID
func UpdateMajorMeaning(db *sql.DB, id int64, input MeaningMajorInput) (*MeaningMajor, error) {
	const query = `
	UPDATE meaning_major
	SET number = $1, position = $2, source = $3, meaning = $4
	WHERE id = $5`
	res, err := db.Exec(query, input.Number, input.Position, input.Source, input.Meaning, id)
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

	updated := &MeaningMajor{
		ID:       id,
		Number:   input.Number,
		Position: input.Position,
		Source:   input.Source,
		Meaning:  input.Meaning,
	}
	return updated, nil
}

// DeleteMajorMeaning deletes a record from the meaning_major table
func DeleteMajorMeaning(db *sql.DB, id int64) error {
	const query = `DELETE FROM meaning_major WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
