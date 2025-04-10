package handlers

import (
	"net/http"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListMajorCardsHandler returns all Major Arcana cards for a given deck
// @Summary List Major Arcana cards
// @Description Returns all major arcana cards from the specified deck
// @Tags cards
// @Accept json
// @Produce json
// @Param deckId query int true "Deck ID (required)"
// @Success 200 {array} models.CardMajor
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /cards/major [get]
func ListMajorCardsHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		deckID, err := useIDParam(c, "deckId")
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}

		cards, err := models.ListMajorCards(a.DB, deckID)
		if err != nil {
			return useHandleDBError(c, err)
		}
		return c.JSON(http.StatusOK, cards)
	}
}

// GetMajorCardByIDHandler returns a specific card by ID
// @Summary Get card by ID
// @Description Retrieves a card by its ID
// @Tags cards
// @Produce json
// @Param id path int true "CardMajor ID"
// @Success 200 {object} models.CardMajor
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /cards/major/{id} [get]
func GetMajorCardByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetMajorCardByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Major Card not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateMajorCardHandler creates a new card
// @Summary Create a new card
// @Description Adds a new card of interpretations
// @Tags cards
// @Accept json
// @Produce json
// @Param card body models.CardMajorInput true "CardMajor data"
// @Success 201 {object} models.IDOnly "Returns the ID of the created card"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /cards/major [post]
func CreateMajorCardHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.CardMajorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateMajorCard(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateMajorCardHandler updates an existing card
// @Summary Update a card
// @Description Updates an existing card
// @Tags cards
// @Accept json
// @Produce json
// @Param id path int true "CardMajor ID"
// @Param card body models.CardMajorInput true "Updated card"
// @Success 200 {object} models.CardMajor
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /cards/major/{id} [put]
func UpdateMajorCardHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.CardMajorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		updated, err := models.UpdateMajorCard(a.DB, id, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Major Card not found")
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteMajorCardHandler deletes a card by ID
// @Summary Delete a card
// @Description Deletes a card by ID
// @Tags cards
// @Param id path int true "CardMajor ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /cards/major/{id} [delete]
func DeleteMajorCardHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteMajorCard(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
