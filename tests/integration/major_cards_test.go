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

func Test_GET_major_cards_returns_all_major_cards(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cards/major?deckId=1", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var cards []models.CardMajor
	err := json.Unmarshal(rec.Body.Bytes(), &cards)
	require.NoError(t, err)
	assert.NotEmpty(t, cards, "expected at least one major card in the list")
}

func Test_POST__major_cards_creates_new_Card(t *testing.T) {
	payload := models.CardMajorInput{
		Name:    "Crocodile 1",
		DeckID:  1,
		OrgName: "Corcodrillo Uno",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cards/major", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.ID > 0, "Expected a valid card ID")
}

func Test_PUT_major_cards_updates_existing_Card(t *testing.T) {
	payload := models.CardMajorInput{
		Name:    "Crocodile 2",
		DeckID:  1,
		OrgName: "Corcodrillo Duo",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cards/major", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	// Get ID
	CardID := getMajorCardId(t, rec)

	// Update
	updated := models.CardMajorInput{
		Name:    "Updated Crocodile",
		DeckID:  1,
		OrgName: "Corcodrillo Duo",
	}
	updatedBody, _ := json.Marshal(updated)
	putReq := httptest.NewRequest(http.MethodPut, "/cards/major/"+strconv.Itoa(CardID), bytes.NewReader(updatedBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(putRec, putReq)

	assert.Equal(t, http.StatusOK, putRec.Code)

	// Parse response and check updated name
	var updatedCard models.Card
	err := json.Unmarshal(putRec.Body.Bytes(), &updatedCard)
	require.NoError(t, err)
	assert.Equal(t, "Updated Crocodile", updatedCard.Name)
	assert.Equal(t, int64(CardID), updatedCard.ID)
}

func Test_DELETE_MajorCard_removes_existing_Card(t *testing.T) {
	payload := models.CardMajorInput{
		Name:    "Crocodile 3",
		DeckID:  1,
		OrgName: "Corcodrillo Tre",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/cards/major", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	CardID := getMajorCardId(t, rec)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/cards/major/"+strconv.Itoa(CardID), nil)
	delRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(delRec, delReq)

	assert.Equal(t, http.StatusNoContent, delRec.Code)
}
