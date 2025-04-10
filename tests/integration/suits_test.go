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

func Test_GET_suits_returns_all_suits(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/suits", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var suits []models.Suit
	err := json.Unmarshal(rec.Body.Bytes(), &suits)
	require.NoError(t, err)
	assert.NotEmpty(t, suits, "expected at least one suit in the list")
}

func Test_POST__suits_creates_new_Suit(t *testing.T) {
	payload := models.SuitInput{
		Name:        "Suit 1",
		Genitive:    "of suit 1",
		Description: "Lorem Ipsum",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/suits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.ID > 0, "Expected a valid suit ID")
}

func Test_PUT_suits_updates_existing_Suit(t *testing.T) {
	payload := models.SuitInput{
		Name:        "Suit 2",
		Genitive:    "of suit 2",
		Description: "Lorem Ipsum",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/suits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	// Get ID
	SuitID := getSuitId(t, rec)

	// Update
	updated := models.SuitInput{
		Name:        "Updated Suit",
		Genitive:    "of updated suit",
		Description: "Lorem Ipsum",
	}
	updatedBody, _ := json.Marshal(updated)
	putReq := httptest.NewRequest(http.MethodPut, "/suits/"+strconv.Itoa(SuitID), bytes.NewReader(updatedBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(putRec, putReq)

	assert.Equal(t, http.StatusOK, putRec.Code)

	// Parse response and check updated name
	var updatedSuit models.Suit
	err := json.Unmarshal(putRec.Body.Bytes(), &updatedSuit)
	require.NoError(t, err)
	assert.Equal(t, "Updated Suit", updatedSuit.Name)
	assert.Equal(t, int64(SuitID), updatedSuit.ID)
}

func getSuitId(t *testing.T, rec *httptest.ResponseRecorder) int {
	var createdSuit models.Suit
	err := json.Unmarshal(rec.Body.Bytes(), &createdSuit)
	require.NoError(t, err)

	id := int(createdSuit.ID)
	return id
}

func Test_DELETE_Suit_removes_existing_Suit(t *testing.T) {
	payload := models.SuitInput{
		Name:        "Suit 3",
		Genitive:    "of suit 3",
		Description: "Lorem Ipsum",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/suits", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	SuitID := getSuitId(t, rec)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/suits/"+strconv.Itoa(SuitID), nil)
	delRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(delRec, delReq)

	assert.Equal(t, http.StatusNoContent, delRec.Code)
}
