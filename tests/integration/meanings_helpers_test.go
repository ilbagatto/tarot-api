package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/stretchr/testify/require"
)

// createTestSource creates a temporary Source for use in tests
func createTestSource(t *testing.T) int64 {
	sourceID, err := createSource()
	require.NoError(t, err)

	t.Cleanup(func() {
		err := deleteSource(*sourceID)
		require.NoError(t, err, "failed to delete test source")
	})

	return *sourceID
}

// createTestMajorMeaning creates a temporary MeaningMajor for a given source
func createTestMajorMeaning(t *testing.T, input models.MeaningMajorInput) int {
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/meanings/major", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	require.True(t, result.ID > 0)

	return int(result.ID)
}

// createTestMinorMeaning creates a temporary MeaningMinor for a given source
func createTestMinorMeaning(t *testing.T, input models.MeaningMinorInput) int {
	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/meanings/minor", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var result models.IDOnly
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	require.True(t, result.ID > 0)

	return int(result.ID)
}
