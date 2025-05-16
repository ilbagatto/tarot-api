package models

import (
	"database/sql"

	"github.com/ilbagatto/tarot-api/internal/utils"
)

// DeckInput represents input format for creating/updating decks
// Please, note that Image field here must be a relative path, unlike
// Image in the other structures for reading.
type DeckInput struct {
	Name        string   `json:"name"`
	Image       string   `json:"image"` // Relative URL
	Description string   `json:"description"`
	Sources     []IDOnly `json:"sources"`
}

// Deck represents a Tarot deck
type Deck struct {
	ID            int64    `json:"id"`
	Name          string   `json:"name"`
	Image         string   `json:"image"` // Full URL
	Description   string   `json:"description"`
	Sources       []Source `json:"sources,omitempty"`
	HasMinorCards bool     `json:"hasMinorCards"`
}

// DeckListItem represents a Tarot deck as list item, without related sources
type DeckListItem struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image"` // Full URL
	Description   string `json:"description"`
	HasMinorCards bool   `json:"hasMinorCards"`
}

// DeckRef is a lightweight deck reference for embedding in Source
type DeckRef struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"` // Full URL
	Description string `json:"description,omitempty"`
}

// fetchDecksFromRows handles rows and returns decks list.
func fetchDecksFromRows(rows *sql.Rows) ([]DeckListItem, error) {
	defer rows.Close()
	var decks []DeckListItem

	for rows.Next() {
		var deck DeckListItem
		if err := rows.Scan(&deck.ID, &deck.Name, &deck.Image, &deck.HasMinorCards, &deck.Description); err != nil {
			return nil, err
		}
		deck.Image = *utils.GetImageURL(deck.Image)
		decks = append(decks, deck)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return decks, nil
}

// ListDecks retrieves all decks
func ListDecks(db *sql.DB) ([]DeckListItem, error) {
	rows, err := db.Query("SELECT id, name, image, has_minor_cards, description FROM deck_with_stats")
	if err != nil {
		return nil, err
	}
	return fetchDecksFromRows(rows)
}

// GetDeckByID retrieves a single deck and its sources
func GetDeckByID(db *sql.DB, deckId int64) (*Deck, error) {
	var deck Deck

	// Main deck query
	row := db.QueryRow("SELECT id, name, image, has_minor_cards, description FROM deck_with_stats WHERE id = $1", deckId)
	if err := row.Scan(&deck.ID, &deck.Name, &deck.Image, &deck.HasMinorCards, &deck.Description); err != nil {
		return nil, err
	}

	deck.Image = *utils.GetImageURL(deck.Image)

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
	query := "INSERT INTO deck (name, image, description) VALUES ($1, $2, $3) RETURNING id"

	var id int64
	if err := db.QueryRow(query, deck.Name, deck.Image, deck.Description).Scan(&id); err != nil {
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
	updateQuery := "UPDATE deck SET name = $1, image = $2, description = $3 WHERE id = $4"
	res, err := tx.Exec(updateQuery, input.Name, input.Image, input.Description, deckID)
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
		Image:       input.Image,
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
