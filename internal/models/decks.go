package models

import (
	"database/sql"
)

type DeckInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Sources     []IDOnly `json:"sources"`
}

// Deck represents a Tarot deck
type Deck struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Sources     []Source `json:"sources,omitempty"`
}

// DeckListItem represents a Tarot deck as list item, without related sources
type DeckListItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// DeckRef is a lightweight deck reference for embedding in Source
type DeckRef struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ListDecks retrieves all decks
func ListDecks(db *sql.DB) ([]DeckListItem, error) {
	var decks []DeckListItem
	rows, err := db.Query("SELECT id, name, description FROM deck")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deck DeckListItem
		if err := rows.Scan(&deck.ID, &deck.Name, &deck.Description); err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return decks, nil
}

// ListNonEmptyDecks retrieves only decks that contain cards and images
func ListNonEmptyDecks(db *sql.DB) ([]DeckListItem, error) {
	var decks []DeckListItem
	rows, err := db.Query("SELECT id, name, description FROM nonempty_decks_view")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deck DeckListItem
		if err := rows.Scan(&deck.ID, &deck.Name, &deck.Description); err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return decks, nil
}

// GetDeckByID retrieves a single deck and its sources
func GetDeckByID(db *sql.DB, deckId int64) (*Deck, error) {
	var deck Deck

	// Main deck query
	row := db.QueryRow("SELECT id, name, description FROM deck WHERE id = $1", deckId)
	if err := row.Scan(&deck.ID, &deck.Name, &deck.Description); err != nil {
		return nil, err
	}

	// Load related sources
	query := `
		SELECT s.id, s.name
		FROM source s
		INNER JOIN deck_source ds ON ds.source = s.id
		WHERE ds.deck = $1
	`
	rows, err := db.Query(query, deck.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var src Source
		if err := rows.Scan(&src.ID, &src.Name); err != nil {
			return nil, err
		}
		deck.Sources = append(deck.Sources, src)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &deck, nil
}

// CreateDeck inserts a new deck into the database and returns the new ID
func CreateDeck(db *sql.DB, deck DeckInput) (*int64, error) {
	query := "INSERT INTO deck (name, description) VALUES ($1, $2) RETURNING id"

	var id int64
	if err := db.QueryRow(query, deck.Name, deck.Description).Scan(&id); err != nil {
		return nil, err
	}

	return &id, nil
}

// UpdateDeck updates an existing deck and its associated sources
func UpdateDeck(db *sql.DB, deckID int64, input DeckInput) (*Deck, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update deck fields
	updateQuery := "UPDATE deck SET name = $1, description = $2 WHERE id = $3"
	res, err := tx.Exec(updateQuery, input.Name, input.Description, deckID)
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

	// Remove existing sources
	if _, err := tx.Exec("DELETE FROM deck_source WHERE deck = $1", deckID); err != nil {
		return nil, err
	}

	// Insert new sources
	for _, src := range input.Sources {
		_, err := tx.Exec("INSERT INTO deck_source (deck, source) VALUES ($1, $2)", deckID, src.ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Build the full deck response with attached sources
	updated := &Deck{
		ID:          deckID,
		Name:        input.Name,
		Description: input.Description,
		Sources:     make([]Source, len(input.Sources)),
	}
	for i, src := range input.Sources {
		updated.Sources[i] = Source{ID: src.ID}
	}

	return updated, nil
}

// DeleteDeck removes a deck by ID
func DeleteDeck(db *sql.DB, deckId int64) error {
	query := "DELETE FROM deck WHERE id = $1"
	_, err := db.Exec(query, deckId)
	return err
}
