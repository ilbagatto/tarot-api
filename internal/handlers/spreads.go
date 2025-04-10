package handlers

import (
	"net/http"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListSpreadsHandler returns all spreads
// @Summary Get all spreads
// @Description Retrieves a list of all available interpretation spreads
// @Tags spreads
// @Produce json
// @Success 200 {array} models.Spread
// @Failure 500 {object} APIResponse
// @Router /spreads [get]
func ListSpreadsHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		spreads, err := models.ListSpreads(a.DB)
		if err != nil {
			return useHandleDBError(c, err)
		}
		return c.JSON(http.StatusOK, spreads)
	}
}

// GetSpreadByIDHandler returns a specific spread by ID
// @Summary Get spread by ID
// @Description Retrieves a spread by its ID
// @Tags spreads
// @Produce json
// @Param id path int true "Spread ID"
// @Success 200 {object} models.Spread
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /spreads/{id} [get]
func GetSpreadByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetSpreadByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Spread not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateSpreadHandler creates a new spread
// @Summary Create a new spread
// @Description Adds a new spread of interpretations
// @Tags spreads
// @Accept json
// @Produce json
// @Param spread body models.SpreadInput true "Spread data"
// @Success 201 {object} models.IDOnly "Returns the ID of the created spread"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /spreads [post]
func CreateSpreadHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.SpreadInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateSpread(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateSpreadHandler updates an existing spread
// @Summary Update a spread
// @Description Updates an existing spread
// @Tags spreads
// @Accept json
// @Produce json
// @Param id path int true "Spread ID"
// @Param spread body models.SpreadInput true "Updated spread"
// @Success 200 {object} models.Spread
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /spreads/{id} [put]
func UpdateSpreadHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.SpreadInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		updated, err := models.UpdateSpread(a.DB, id, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Spread not found")
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteSpreadHandler deletes a spread by ID
// @Summary Delete a spread
// @Description Deletes a spread by ID
// @Tags spreads
// @Param id path int true "Spread ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /spreads/{id} [delete]
func DeleteSpreadHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteSpread(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
