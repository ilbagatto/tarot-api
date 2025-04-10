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

func Test_GET__sources_returns_all_sources(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/sources", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var sources []models.Source
	err := json.Unmarshal(rec.Body.Bytes(), &sources)
	require.NoError(t, err)
	assert.NotEmpty(t, sources, "expected at least one source in the list")
}

func Test_POST__sources_creates_new_Source(t *testing.T) {
	payload := models.SourceInput{
		Name:  "Test Source",
		Decks: []models.IDOnly{{ID: 1}},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.ID > 0, "Expected a valid suit ID")
}

func Test_PUT__sources_updates_existing_Source(t *testing.T) {
	payload := models.SourceInput{
		Name:  "Another Test Source",
		Decks: []models.IDOnly{{ID: 1}},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	// Get ID
	SourceID := getSourceId(t, rec)

	// Update
	updated := models.SourceInput{
		Name:  "Updated Source",
		Decks: []models.IDOnly{{ID: 1}},
	}
	updatedBody, _ := json.Marshal(updated)
	putReq := httptest.NewRequest(http.MethodPut, "/sources/"+strconv.Itoa(SourceID), bytes.NewReader(updatedBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(putRec, putReq)

	assert.Equal(t, http.StatusOK, putRec.Code)
	// Parse response and check updated name
	var updatedSource models.Source
	err := json.Unmarshal(putRec.Body.Bytes(), &updatedSource)
	require.NoError(t, err)
	assert.Equal(t, "Updated Source", updatedSource.Name)
	assert.Equal(t, int64(SourceID), updatedSource.ID)

}

func getSourceId(t *testing.T, rec *httptest.ResponseRecorder) int {
	var createdSource models.Source
	err := json.Unmarshal(rec.Body.Bytes(), &createdSource)
	require.NoError(t, err)

	SourceID := int(createdSource.ID)
	return SourceID
}

func Test_DELETE__Source_removes_existing_Source(t *testing.T) {
	payload := models.SourceInput{
		Name:  "One more Test Source",
		Decks: []models.IDOnly{{ID: 1}},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/sources", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	SourceID := getSourceId(t, rec)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/sources/"+strconv.Itoa(SourceID), nil)
	delRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(delRec, delReq)

	assert.Equal(t, http.StatusNoContent, delRec.Code)
}
