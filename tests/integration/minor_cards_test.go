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

func Test_GET_minor_cards_returns_all_minor_cards(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cards/minor?deckId=1", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var cards []models.CardMinor
	err := json.Unmarshal(rec.Body.Bytes(), &cards)
	require.NoError(t, err)
	assert.NotEmpty(t, cards, "expected at least one minor card in the list")
}

func Test_POST__minor_cards_creates_new_Card(t *testing.T) {
	deckID, err := createDeck()
	require.NoError(t, err)
	defer func() {
		err := deleteDeck(*deckID)
		require.NoError(t, err, "failed to delete test deck")
	}()
	payload := models.CardMinorInput{
		DeckID: *deckID,
		SuitID: 1,
		RankID: 1,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cards/minor", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var result models.IDOnly
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.ID > 0, "Expected a valid card ID")
}

func Test_PUT_minor_cards_updates_existing_Card(t *testing.T) {
	deckID, err := createDeck()
	require.NoError(t, err)
	defer func() {
		err := deleteDeck(*deckID)
		require.NoError(t, err, "failed to delete test deck")
	}()

	payload := models.CardMinorInput{
		DeckID: *deckID,
		SuitID: 1,
		RankID: 1,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cards/minor", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	// Get ID
	CardID := getMinorCardId(t, rec)

	// Update
	updated := models.CardMinorInput{
		DeckID: *deckID,
		SuitID: 2,
		RankID: 1,
	}
	updatedBody, _ := json.Marshal(updated)
	putReq := httptest.NewRequest(http.MethodPut, "/cards/minor/"+strconv.Itoa(CardID), bytes.NewReader(updatedBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(putRec, putReq)

	assert.Equal(t, http.StatusNoContent, putRec.Code)
}

func Test_DELETE_MinorCard_removes_existing_Card(t *testing.T) {
	deckID, err := createDeck()
	require.NoError(t, err)
	defer func() {
		err := deleteDeck(*deckID)
		require.NoError(t, err, "failed to delete test deck")
	}()

	payload := models.CardMinorInput{
		DeckID: *deckID,
		SuitID: 1,
		RankID: 1,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cards/minor", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	CardID := getMajorCardId(t, rec)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/cards/minor/"+strconv.Itoa(CardID), nil)
	delRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(delRec, delReq)

	assert.Equal(t, http.StatusNoContent, delRec.Code)
}
