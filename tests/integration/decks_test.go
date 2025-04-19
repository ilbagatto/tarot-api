package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GET__decks_returns_all_decks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/decks", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var decks []models.Deck
	err := json.Unmarshal(rec.Body.Bytes(), &decks)
	require.NoError(t, err)
	assert.NotEmpty(t, decks, "expected at least one deck in the list")
}

func Test_GET__decks_with_hasCards_param(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/decks?hasCards=yes", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var decks []models.Deck
	err := json.Unmarshal(rec.Body.Bytes(), &decks)
	require.NoError(t, err)
	assert.NotEmpty(t, decks, "expected at least one deck in the list")
}

func Test_POST__decks_creates_new_deck(t *testing.T) {
	payload := models.DeckInput{
		Name:        "Test Deck",
		Image:       "image.png",
		Description: "Test description",
		Sources:     []models.IDOnly{{ID: 1}},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/decks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.ID > 0, "Expected a valid deck ID")
}

// Test_POST__decks_with_duplicate_name_returns_409 checks that duplicate names are rejected
func Test_POST__decks_with_duplicate_name_returns_409(t *testing.T) {
	payload := models.DeckInput{
		Name:        "Duplicate Deck",
		Image:       "image.png",
		Description: "First instance",
		Sources:     []models.IDOnly{{ID: 1}},
	}
	body, _ := json.Marshal(payload)

	// Create first deck
	req1 := httptest.NewRequest(http.MethodPost, "/decks", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec1, req1)
	require.Equal(t, http.StatusCreated, rec1.Code)

	// Create duplicate
	req2 := httptest.NewRequest(http.MethodPost, "/decks", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec2, req2)

	assert.Equal(t, http.StatusConflict, rec2.Code)
	assert.Contains(t, rec2.Body.String(), "Duplicate entry")
}

// Test_GET__nonexistent_deck_returns_404 checks GET for invalid ID
func Test_GET__nonexistent_deck_returns_404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/decks/999999", nil)
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// Test_PUT__nonexistent_deck_returns_404 checks update for nonexistent deck
func Test_PUT__nonexistent_deck_returns_404(t *testing.T) {
	payload := models.DeckInput{
		Name:        "Nonexistent Deck",
		Description: "Should not exist",
		Sources:     []models.IDOnly{{ID: 1}},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/decks/999999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func Test_PUT__decks_updates_existing_deck(t *testing.T) {
	payload := models.DeckInput{
		Name:        "Deck to Update",
		Image:       "image.png",
		Description: "Old description",
		Sources:     []models.IDOnly{{ID: 1}},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/decks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	// Get ID
	deckID := getDeckId(t, rec)

	// Update
	updated := models.DeckInput{
		Name:        "Updated Deck",
		Image:       "image.png",
		Description: "Updated description",
		Sources:     []models.IDOnly{{ID: 1}},
	}
	updatedBody, _ := json.Marshal(updated)
	putReq := httptest.NewRequest(http.MethodPut, "/decks/"+strconv.Itoa(deckID), bytes.NewReader(updatedBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(putRec, putReq)

	assert.Equal(t, http.StatusOK, putRec.Code)

	// Parse response and check updated name
	var updatedDeck models.Deck
	err := json.Unmarshal(putRec.Body.Bytes(), &updatedDeck)
	require.NoError(t, err)
	assert.Equal(t, "Updated Deck", updatedDeck.Name)
	assert.Equal(t, int64(deckID), updatedDeck.ID)

}

func getDeckId(t *testing.T, rec *httptest.ResponseRecorder) int {
	var createdDeck models.Deck
	err := json.Unmarshal(rec.Body.Bytes(), &createdDeck)
	require.NoError(t, err)
	deckID := int(createdDeck.ID)
	return deckID
}

func Test_DELETE__deck_removes_existing_deck(t *testing.T) {
	// Create a deck
	payload := models.DeckInput{
		Name:        "Deck to Delete",
		Image:       "image.png",
		Description: "To be deleted",
		Sources:     []models.IDOnly{{ID: 1}},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/decks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	deckID := getDeckId(t, rec)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/decks/"+strconv.Itoa(deckID), nil)
	delRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(delRec, delReq)

	assert.Equal(t, http.StatusNoContent, delRec.Code)
}

// Test_DELETE__nonexistent_deck_returns_404 checks DELETE for nonexistent deck
func Test_DELETE__nonexistent_deck_returns_204(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/decks/999999", nil)
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
