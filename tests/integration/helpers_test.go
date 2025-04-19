package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/ilbagatto/tarot-api/internal/testutils"
	"github.com/stretchr/testify/require"
)

func createDeck() (*int64, error) {
	deck := models.DeckInput{
		Name:        testutils.RandomString(10, 50),
		Image:       "image.png",
		Description: "TODO",
		Sources:     []models.IDOnly{{ID: 1}},
	}
	body, _ := json.Marshal(deck)
	req := httptest.NewRequest(http.MethodPost, "/decks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	var createdItem models.Deck
	if err := json.Unmarshal(rec.Body.Bytes(), &createdItem); err != nil {
		return nil, err
	}
	return &createdItem.ID, nil
}

func deleteDeck(deckID int64) error {
	req := httptest.NewRequest(http.MethodDelete, "/decks/"+strconv.FormatInt(deckID, 10), nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	if rec.Code == http.StatusNoContent {
		return nil // everything is fine
	}
	return fmt.Errorf("failed to delete deck %d: HTTP %d – %s", deckID, rec.Code, rec.Body.String())
}

func getMinorCardId(t *testing.T, rec *httptest.ResponseRecorder) int {
	var card models.CardMinor
	err := json.Unmarshal(rec.Body.Bytes(), &card)
	require.NoError(t, err)

	id := int(card.ID)
	return id
}

func getMajorCardId(t *testing.T, rec *httptest.ResponseRecorder) int {
	var card models.CardMajor
	err := json.Unmarshal(rec.Body.Bytes(), &card)
	require.NoError(t, err)

	id := int(card.ID)
	return id
}

func createSource() (*int64, error) {
	payload := models.SourceInput{
		Name:  testutils.RandomString(10, 50),
		Decks: []models.IDOnly{{ID: 1}},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	var createdItem models.Source
	if err := json.Unmarshal(rec.Body.Bytes(), &createdItem); err != nil {
		return nil, err
	}
	return &createdItem.ID, nil
}

func deleteSource(sourceID int64) error {
	req := httptest.NewRequest(http.MethodDelete, "/sources/"+strconv.FormatInt(sourceID, 10), nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	if rec.Code == http.StatusNoContent {
		return nil // everything is fine
	}
	return fmt.Errorf("failed to delete source %d: HTTP %d – %s", sourceID, rec.Code, rec.Body.String())
}
