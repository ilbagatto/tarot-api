package models

import (
	"database/sql"
)

// Source represents a source of card interpretations
type Source struct {
	ID    int64     `json:"id"`
	Name  string    `json:"name"`
	Decks []DeckRef `json:"decks,omitempty"`
}

// SourceInput represents input format for creating/updating sources
type SourceInput struct {
	Name  string   `json:"name" example:"Мишель Моран"`
	Decks []IDOnly `json:"decks"`
}

// SourceListItem represents a source without related decks, as represented in list.
type SourceListItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name" example:"Мишель Моран"`
}

// ListSources retrieves all sources
func ListSources(db *sql.DB) ([]SourceListItem, error) {
	var sources []SourceListItem
	rows, err := db.Query("SELECT id, name FROM source")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var src SourceListItem
		if err := rows.Scan(&src.ID, &src.Name); err != nil {
			return nil, err
		}
		sources = append(sources, src)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sources, nil
}

// GetSourceByID retrieves a single source by ID
func GetSourceByID(db *sql.DB, id int64) (*Source, error) {
	var src Source

	// Fetch the source
	row := db.QueryRow("SELECT id, name FROM source WHERE id = $1", id)
	if err := row.Scan(&src.ID, &src.Name); err != nil {
		return nil, err
	}

	// Fetch related decks
	rows, err := db.Query(`
		SELECT d.id, d.name, d.description
		FROM deck d
		INNER JOIN deck_source ds ON ds.deck = d.id
		WHERE ds.source = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deck DeckRef
		if err := rows.Scan(&deck.ID, &deck.Name, &deck.Description); err != nil {
			return nil, err
		}
		src.Decks = append(src.Decks, deck)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &src, nil
}

// CreateSource inserts a new source
func CreateSource(db *sql.DB, input SourceInput) (*int64, error) {
	query := "INSERT INTO source (name) VALUES ($1) RETURNING id"
	var id int64
	err := db.QueryRow(query, input.Name).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// UpdateSource updates an existing source
func UpdateSource(db *sql.DB, sourceID int64, input SourceInput) (*Source, error) {
	query := "UPDATE source SET name = $1 WHERE id = $2"
	res, err := db.Exec(query, input.Name, sourceID)
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

	updated := &Source{
		ID:   sourceID,
		Name: input.Name,
	}

	return updated, nil
}

// DeleteSource deletes a source by ID
func DeleteSource(db *sql.DB, id int64) error {
	query := "DELETE FROM source WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}
