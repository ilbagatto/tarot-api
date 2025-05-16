package handlers

import (
	"net/http"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListDecksHandler returns a list of all Tarot decks
// @Summary Get all decks
// @Description Retrieves a list of available Tarot decks. Optionally filters decks that contain cards.
// @Tags decks
// @Produce json
// @Success 200 {array} models.DeckListItem
// @Failure 500 {object} map[string]string
// @Router /decks [get]
func ListDecksHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var decks []models.DeckListItem
		var err error
		decks, err = models.ListDecks(a.DB)

		if err != nil {
			return useHandleDBError(c, err)
		}
		return c.JSON(http.StatusOK, decks)
	}
}

// GetDeckByIDHandler returns a deck by ID, including its associated sources
// @Summary Get deck by ID
// @Description Retrieves a deck with its associated interpretation sources
// @Tags decks
// @Produce json
// @Param id path int true "Deck ID"
// @Success 200 {object} models.Deck
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /decks/{id} [get]
func GetDeckByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		deckID, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}

		deck, err := models.GetDeckByID(a.DB, deckID)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Deck not found")
		}

		return c.JSON(http.StatusOK, deck)
	}
}

// CreateDeckHandler creates a new deck
// @Summary Create a new deck
// @Description Adds a new Tarot deck and links to existing sources (by ID)
// @Tags decks
// @Accept json
// @Produce json
// @Param deck body models.DeckInput true "Deck input"
// @Success 201 {object} models.IDOnly
// @Failure 400 {object} APIResponse
// @Failure 409 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /decks [post]
func CreateDeckHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.DeckInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}

		id, err := models.CreateDeck(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateDeckHandler updates an existing deck
// @Summary Update a deck
// @Description Updates the specified Tarot deck and its sources
// @Tags decks
// @Accept json
// @Produce json
// @Param id path int true "Deck ID"
// @Param deck body models.DeckInput true "Deck input"
// @Success 200 {object} models.Deck
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 409 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /decks/{id} [put]
func UpdateDeckHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		deckID, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}

		var input models.DeckInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}

		deck, err := models.UpdateDeck(a.DB, deckID, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Deck not found")
		}

		return c.JSON(http.StatusOK, deck)
	}
}

// DeleteDeckHandler handles DELETE /decks/:id
// @Summary Delete deck by ID
// @Description Deletes a deck and its associations
// @Tags decks
// @Produce json
// @Param id path int true "Deck ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /decks/{id} [delete]
func DeleteDeckHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		deckID, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}

		if err := models.DeleteDeck(a.DB, deckID); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
