package handlers

import (
	"net/http"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListMinorCardsHandler returns all Minor Arcana cards for a given deck
// @Summary List Minor Arcana cards
// @Description Returns all minor arcana cards from the specified deck
// @Tags cards
// @Accept json
// @Produce json
// @Param deckId query int true "Deck ID (required)"
// @Success 200 {array} models.CardMinor
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /cards/minor [get]
func ListMinorCardsHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		deckID, err := useIDParam(c, "deckId")
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}

		cards, err := models.ListMinorCards(a.DB, deckID)
		if err != nil {
			return useHandleDBError(c, err)
		}
		return c.JSON(http.StatusOK, cards)
	}
}

// GetMinorCardByIDHandler returns a specific card by ID
// @Summary Get card by ID
// @Description Retrieves a card by its ID
// @Tags cards
// @Produce json
// @Param id path int true "CardMinor ID"
// @Success 200 {object} models.CardMinor
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /cards/minor/{id} [get]
func GetMinorCardByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetMinorCardByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Minor Card not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateMinorCardHandler creates a new card
// @Summary Create a new card
// @Description Adds a new card of interpretations
// @Tags cards
// @Accept json
// @Produce json
// @Param card body models.CardMinorInput true "CardMinor data"
// @Success 201 {object} models.IDOnly "Returns the ID of the created card"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /cards/minor [post]
func CreateMinorCardHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.CardMinorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateMinorCard(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateMinorCardHandler updates an existing card
// @Summary Update a card
// @Description Updates an existing card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path int true "CardMinor ID"
// @Param card body models.CardMinorInput true "Updated card"
// @Success 200 {object} models.CardMinor
// @Success 204 "No Content"
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /cards/minor/{id} [put]
func UpdateMinorCardHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.CardMinorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.UpdateMinorCard(a.DB, id, input); err != nil {
			return useHandleNotFoundOrDBError(c, err, "Minor Card not found")
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// DeleteMinorCardHandler deletes a card by ID
// @Summary Delete a card
// @Description Deletes a card by ID
// @Tags cards
// @Param id path int true "CardMinor ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /cards/minor/{id} [delete]
func DeleteMinorCardHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteMinorCard(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
