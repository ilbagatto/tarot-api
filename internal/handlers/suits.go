package handlers

import (
	"net/http"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListSuitsHandler returns all suits
// @Summary Get all suits
// @Description Retrieves a list of all available interpretation suits
// @Tags suits
// @Produce json
// @Success 200 {array} models.Suit
// @Failure 500 {object} APIResponse
// @Router /suits [get]
func ListSuitsHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		suits, err := models.ListSuits(a.DB)
		if err != nil {
			return useHandleDBError(c, err)
		}
		return c.JSON(http.StatusOK, suits)
	}
}

// GetSuitByIDHandler returns a specific suit by ID
// @Summary Get suit by ID
// @Description Retrieves a suit by its ID
// @Tags suits
// @Produce json
// @Param id path int true "Suit ID"
// @Success 200 {object} models.Suit
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /suits/{id} [get]
func GetSuitByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetSuitByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Suit not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateSuitHandler creates a new suit
// @Summary Create a new suit
// @Description Adds a new suit of interpretations
// @Tags suits
// @Accept json
// @Produce json
// @Param suit body models.SuitInput true "Suit data"
// @Success 201 {object} models.IDOnly "Returns the ID of the created suit"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /suits [post]
func CreateSuitHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.SuitInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateSuit(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateSuitHandler updates an existing suit
// @Summary Update a suit
// @Description Updates an existing suit
// @Tags suits
// @Accept json
// @Produce json
// @Param id path int true "Suit ID"
// @Param suit body models.SuitInput true "Updated suit"
// @Success 200 {object} models.Suit
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /suits/{id} [put]
func UpdateSuitHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.SuitInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		updated, err := models.UpdateSuit(a.DB, id, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Suit not found")
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteSuitHandler deletes a suit by ID
// @Summary Delete a suit
// @Description Deletes a suit by ID
// @Tags suits
// @Param id path int true "Suit ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /suits/{id} [delete]
func DeleteSuitHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteSuit(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
