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

func Test_GET_ranks_returns_all_ranks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ranks", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var ranks []models.Rank
	err := json.Unmarshal(rec.Body.Bytes(), &ranks)
	require.NoError(t, err)
	assert.NotEmpty(t, ranks, "expected at least one rank in the list")
}

func Test_POST__ranks_creates_new_Rank(t *testing.T) {
	payload := models.RankInput{
		Name: "Rank 1",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/ranks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)

	assert.True(t, result.ID > 0, "Expected a valid suit ID")
}

func Test_PUT_ranks_updates_existing_Rank(t *testing.T) {
	payload := models.RankInput{
		Name: "Rank 2",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/ranks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	// Get ID
	RankID := getRankId(t, rec)

	// Update
	updated := models.RankInput{
		Name: "Updated Rank",
	}
	updatedBody, _ := json.Marshal(updated)
	putReq := httptest.NewRequest(http.MethodPut, "/ranks/"+strconv.Itoa(RankID), bytes.NewReader(updatedBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(putRec, putReq)

	assert.Equal(t, http.StatusOK, putRec.Code)

	// Parse response and check updated name
	var updatedRank models.Rank
	err := json.Unmarshal(putRec.Body.Bytes(), &updatedRank)
	require.NoError(t, err)
	assert.Equal(t, "Updated Rank", updatedRank.Name)
	assert.Equal(t, int64(RankID), updatedRank.ID)
}

func getRankId(t *testing.T, rec *httptest.ResponseRecorder) int {
	var createdRank models.Rank
	err := json.Unmarshal(rec.Body.Bytes(), &createdRank)
	require.NoError(t, err)

	RankID := int(createdRank.ID)
	return RankID
}

func Test_DELETE_Rank_removes_existing_Rank(t *testing.T) {
	payload := models.RankInput{
		Name: "Rank 3",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/ranks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	RankID := getRankId(t, rec)

	// Delete
	delReq := httptest.NewRequest(http.MethodDelete, "/ranks/"+strconv.Itoa(RankID), nil)
	delRec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(delRec, delReq)

	assert.Equal(t, http.StatusNoContent, delRec.Code)
}
