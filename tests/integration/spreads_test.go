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

func Test_GET__spreads_returns_all_spreads(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/spreads", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var spreads []models.Spread
	err := json.Unmarshal(rec.Body.Bytes(), &spreads)
	require.NoError(t, err)
	assert.NotEmpty(t, spreads, "expected at least one spread in the list")
}

func Test_POST__spreads_creates_new_Spread(t *testing.T) {
	payload := models.SpreadInput{
		Name:        "A Test Spread 1",
		MajorArcana: true,
		MinorArcana: true,
		UpsideDown:  true,
		NumCards:    3,
		Description: "Lorem ipsum",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/spreads", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.ID > 0, "Expected a valid spread ID")
}

func Test_PUT__spreads_updates_existing_Spread(t *testing.T) {
	payload := models.SpreadInput{
		Name:        "A Test Spread 2",
		MajorArcana: true,
		MinorArcana: true,
		UpsideDown:  true,
		NumCards:    3,
		Description: "Lorem ipsum",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/spreads", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	// Get ID
	SpreadID := getSpreadId(t, rec)

	// Update
	updated := models.SpreadInput{
		Name:        "Updated Spread",
		MajorArcana: true,
		MinorArcana: true,
		UpsideDown:  true,
		NumCards:    3,
		Description: "Lorem ipsum",
	}
	updatedBody, _ := json.Marshal(updated)
	putReq := httptest.NewRequest(http.MethodPut, "/spreads/"+strconv.Itoa(SpreadID), bytes.NewReader(updatedBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(putRec, putReq)

	assert.Equal(t, http.StatusOK, putRec.Code)

	// Parse response and check updated name
	var updatedSpread models.Spread
	err := json.Unmarshal(putRec.Body.Bytes(), &updatedSpread)
	require.NoError(t, err)
	assert.Equal(t, "Updated Spread", updatedSpread.Name)
	assert.Equal(t, int64(SpreadID), updatedSpread.ID)
}

func getSpreadId(t *testing.T, rec *httptest.ResponseRecorder) int {
	var createdSpread models.Spread
	err := json.Unmarshal(rec.Body.Bytes(), &createdSpread)
	require.NoError(t, err)

	SpreadID := int(createdSpread.ID)
	return SpreadID
}

func Test_DELETE__Spread_removes_existing_Spread(t *testing.T) {
	payload := models.SpreadInput{
		Name:        "A Test Spread 3",
		MajorArcana: true,
		MinorArcana: true,
		UpsideDown:  true,
		NumCards:    3,
		Description: "Lorem ipsum",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/spreads", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	SpreadID := getSpreadId(t, rec)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/spreads/"+strconv.Itoa(SpreadID), nil)
	delRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(delRec, delReq)

	assert.Equal(t, http.StatusNoContent, delRec.Code)
}
