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

func Test_POST__minor_meanings_creates_new_MinorMeaning(t *testing.T) {
	sourceID := createTestSource(t)

	id := createTestMinorMeaning(t, models.MeaningMinorInput{
		Suit:     1,
		Rank:     1,
		Position: "straight",
		Source:   sourceID,
		Meaning:  "All will be fine",
	})

	assert.True(t, id > 0)
}

func Test_GET__minor_meanings_returns_filtered_results(t *testing.T) {
	sourceID := createTestSource(t)

	createTestMinorMeaning(t, models.MeaningMinorInput{
		Suit:     1,
		Rank:     1,
		Position: "reverted",
		Source:   sourceID,
		Meaning:  "Hidden fire",
	})

	req := httptest.NewRequest(http.MethodGet, "/meanings/minor?source="+strconv.FormatInt(sourceID, 10)+"&position=reverted", nil)
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Hidden fire")
}

func Test_PUT__minor_meanings_updates_existing_entry(t *testing.T) {
	sourceID := createTestSource(t)

	id := createTestMinorMeaning(t, models.MeaningMinorInput{
		Suit:     1,
		Rank:     1,
		Position: "straight",
		Source:   sourceID,
		Meaning:  "Everything will be broken",
	})

	updated := models.MeaningMinorInput{
		Suit:     1,
		Rank:     1,
		Position: "straight",
		Source:   sourceID,
		Meaning:  "Calm wisdom",
	}

	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPut, "/meanings/minor/"+strconv.Itoa(id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	testApp.App.Echo.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Calm wisdom")
}

func Test_DELETE__minor_meanings_removes_existing_entry(t *testing.T) {
	sourceID := createTestSource(t)

	id := createTestMinorMeaning(t, models.MeaningMinorInput{
		Suit:     1,
		Rank:     1,
		Position: "straight",
		Source:   sourceID,
		Meaning:  "All will be fine",
	})

	req := httptest.NewRequest(http.MethodDelete, "/meanings/minor/"+strconv.Itoa(id), nil)
	rec := httptest.NewRecorder()
	testApp.App.Echo.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
