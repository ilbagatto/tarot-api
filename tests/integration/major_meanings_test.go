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
)

func Test_POST__major_meanings_creates_new_MajorMeaning(t *testing.T) {
	sourceID := createTestSource(t)

	id := createTestMajorMeaning(t, models.MeaningMajorInput{
		Number:   1,
		Position: "straight",
		Source:   sourceID,
		Meaning:  "All will be fine",
	})

	assert.True(t, id > 0)
}

func Test_GET__major_meanings_returns_filtered_results(t *testing.T) {
	sourceID := createTestSource(t)

	createTestMajorMeaning(t, models.MeaningMajorInput{
		Number:   7,
		Position: "reverted",
		Source:   sourceID,
		Meaning:  "Hidden fire",
	})

	req := httptest.NewRequest(http.MethodGet, "/meanings/major?source="+strconv.FormatInt(sourceID, 10)+"&position=reverted", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Hidden fire")
}

func Test_PUT__major_meanings_updates_existing_entry(t *testing.T) {
	sourceID := createTestSource(t)

	id := createTestMajorMeaning(t, models.MeaningMajorInput{
		Number:   12,
		Position: "reverted",
		Source:   sourceID,
		Meaning:  "Past anger",
	})

	updated := models.MeaningMajorInput{
		Number:   12,
		Position: "reverted",
		Source:   sourceID,
		Meaning:  "Calm wisdom",
	}

	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPut, "/meanings/major/"+strconv.Itoa(id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Calm wisdom")
}

func Test_DELETE__major_meanings_removes_existing_entry(t *testing.T) {
	sourceID := createTestSource(t)

	id := createTestMajorMeaning(t, models.MeaningMajorInput{
		Number:   9,
		Position: "straight",
		Source:   sourceID,
		Meaning:  "Silent retreat",
	})

	req := httptest.NewRequest(http.MethodDelete, "/meanings/major/"+strconv.Itoa(id), nil)
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
