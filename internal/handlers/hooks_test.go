package handlers

import (
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_useIDParam_QueryParam(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/?deckId=42", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	id, err := useIDParam(c, "deckId")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), id)
}

func Test_useIDParam_PathParam(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/decks/42", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// manually set path param
	c.SetParamNames("deckId")
	c.SetParamValues("42")

	id, err := useIDParam(c, "deckId")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), id)
}

func Test_useIDParam_DefaultParamName(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/?id=123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	id, err := useIDParam(c) // default is "id"
	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
}

func Test_useIDParam_MissingParam(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_, err := useIDParam(c, "missing")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing parameter")
}

func Test_useIDParam_InvalidParam(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("GET", "/?deckId=abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_, err := useIDParam(c, "deckId")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be an integer")
}
