package models

import (
	"database/sql"
)

type Suit struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Genitive    string `json:"genitive"`
	Description string `json:"description,omitempty"`
}

type SuitInput struct {
	Name        string `json:"name"`
	Genitive    string `json:"genitive"`
	Description string `json:"description,omitempty"`
}

// ListSuits retrieves all suits
func ListSuits(db *sql.DB) ([]Suit, error) {
	rows, err := db.Query(`SELECT id, name, genitive, COALESCE(description, '') FROM suit ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suits []Suit
	for rows.Next() {
		var s Suit
		if err := rows.Scan(&s.ID, &s.Name, &s.Genitive, &s.Description); err != nil {
			return nil, err
		}
		suits = append(suits, s)
	}
	return suits, rows.Err()
}

// GetSuitByID retrieves a single suit by ID
func GetSuitByID(db *sql.DB, id int64) (*Suit, error) {
	var s Suit
	row := db.QueryRow(`SELECT id, name, genitive, COALESCE(description, '') FROM suit WHERE id = $1`, id)
	if err := row.Scan(&s.ID, &s.Name, &s.Genitive, &s.Description); err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSuit inserts a new suit
func CreateSuit(db *sql.DB, s SuitInput) (*int64, error) {
	var id int64
	if err := db.QueryRow(
		`INSERT INTO suit (name, genitive, description) VALUES ($1, $2, $3) RETURNING id`,
		s.Name, s.Genitive, s.Description,
	).Scan(&id); err != nil {
		return nil, err
	}
	return &id, nil
}

// UpdateSuit updates an existing suit
func UpdateSuit(db *sql.DB, suitID int64, s SuitInput) (*Suit, error) {
	res, err := db.Exec(
		`UPDATE suit SET name = $1, genitive = $2, description = $3 WHERE id = $4`,
		s.Name, s.Genitive, s.Description, suitID,
	)
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

	updated := &Suit{
		ID:          suitID,
		Name:        s.Name,
		Genitive:    s.Genitive,
		Description: s.Description,
	}

	return updated, nil
}

// DeleteSuit deletes a suit by ID
func DeleteSuit(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM suit WHERE id = $1`, id)
	return err
}
