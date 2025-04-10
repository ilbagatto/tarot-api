package models

import (
	"database/sql"
)

type SpreadInput struct {
	Name        string `json:"name"`
	MajorArcana bool   `json:"major_arcana"`
	MinorArcana bool   `json:"minor_arcana"`
	UpsideDown  bool   `json:"upside_down"`
	NumCards    int16  `json:"num_cards"`
	Description string `json:"description,omitempty"`
}

type Spread struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	MajorArcana bool   `json:"major_arcana"`
	MinorArcana bool   `json:"minor_arcana"`
	UpsideDown  bool   `json:"upside_down"`
	NumCards    int16  `json:"num_cards"`
	Description string `json:"description,omitempty"`
}

// ListSpreads retrieves all spreads
func ListSpreads(db *sql.DB) ([]Spread, error) {
	rows, err := db.Query(`SELECT id, name, major_arcana, minor_arcana, upside_down, num_cards, description FROM spread`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spreads []Spread
	for rows.Next() {
		var s Spread
		if err := rows.Scan(&s.ID, &s.Name, &s.MajorArcana, &s.MinorArcana, &s.UpsideDown, &s.NumCards, &s.Description); err != nil {
			return nil, err
		}
		spreads = append(spreads, s)
	}
	return spreads, rows.Err()
}

// GetSpreadByID retrieves a single spread by ID
func GetSpreadByID(db *sql.DB, id int64) (*Spread, error) {
	var s Spread
	row := db.QueryRow(`SELECT id, name, major_arcana, minor_arcana, upside_down, num_cards, description FROM spread WHERE id = $1`, id)
	if err := row.Scan(&s.ID, &s.Name, &s.MajorArcana, &s.MinorArcana, &s.UpsideDown, &s.NumCards, &s.Description); err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSpread inserts a new spread
func CreateSpread(db *sql.DB, s SpreadInput) (*int64, error) {
	query := `
	INSERT INTO spread (name, major_arcana, minor_arcana, upside_down, num_cards, description)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
	var id int64
	if err := db.QueryRow(query, s.Name, s.MajorArcana, s.MinorArcana, s.UpsideDown, s.NumCards, s.Description).Scan(&id); err != nil {
		return nil, err
	}
	return &id, nil
}

// UpdateSpread updates an existing spread
func UpdateSpread(db *sql.DB, spreadID int64, s SpreadInput) (*Spread, error) {
	res, err := db.Exec(`
	UPDATE spread
	SET name = $1, major_arcana = $2, minor_arcana = $3, upside_down = $4, num_cards = $5, description = $6
	WHERE id = $7`,
		s.Name, s.MajorArcana, s.MinorArcana, s.UpsideDown, s.NumCards, s.Description, spreadID)
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

	updated := &Spread{
		ID:          spreadID,
		Name:        s.Name,
		MajorArcana: s.MajorArcana,
		MinorArcana: s.MinorArcana,
		UpsideDown:  s.UpsideDown,
		NumCards:    s.NumCards,
		Description: s.Description,
	}

	return updated, err
}

// DeleteSpread deletes a spread by ID
func DeleteSpread(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM spread WHERE id = $1", id)
	return err
}
