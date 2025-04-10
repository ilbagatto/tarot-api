package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListMajorMeaningsHandler returns filtered interpretations for Major Arcana cards
// @Summary Get Major Arcana meanings
// @Description Returns a list of meanings for major arcana cards with optional filters
// @Tags meanings
// @Accept json
// @Produce json
// @Param number query int false "Card number (optional)"
// @Param position query string false "Card position" Enums(straight, reverted)
// @Param source query int false "Source ID (optional)"
// @Success 200 {array} models.MeaningMajor
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /meanings/major [get]
func ListMajorMeaningsHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		filters := map[string]any{}

		if number := c.QueryParam("number"); number != "" {
			if n, err := strconv.Atoi(number); err == nil {
				filters["number"] = n
			} else {
				return SendError(c, http.StatusBadRequest, fmt.Errorf("invalid number"))
			}
		}

		if position := c.QueryParam("position"); position != "" {
			if position := c.QueryParam("position"); position != "" {
				if !validPositions[position] {
					return SendError(c, http.StatusBadRequest, fmt.Errorf("invalid position"))
				}
				filters["position"] = position
			}
		}

		if source := c.QueryParam("source"); source != "" {
			if s, err := strconv.ParseInt(source, 10, 64); err == nil {
				filters["source"] = s
			} else {
				return SendError(c, http.StatusBadRequest, fmt.Errorf("invalid source"))
			}
		}

		result, err := models.ListMajorMeanings(a.DB, filters)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusOK, result)
	}
}

// GetMajorMeaningByIDHandler returns a specific MajorMeaning by ID
// @Summary Get MajorMeaning by ID
// @Description Retrieves a MajorMeaning by its ID
// @Tags meanings
// @Produce json
// @Param id path int true "MajorMeaning ID"
// @Success 200 {object} models.MeaningMajor
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /meanings/major/{id} [get]
func GetMajorMeaningByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetMajorMeaningByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "MajorMeaning not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateMajorMeaningHandler creates a new MajorMeaning
// @Summary Create a new MajorMeaning
// @Description Adds a new MajorMeaning of interpretations
// @Tags meanings
// @Accept json
// @Produce json
// @Param MajorMeaning body models.MeaningMajorInput true "MajorMeaning data"
// @Success 201 {object} models.IDOnly "Returns the ID of the created MajorMeaning"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /meanings/major [post]
func CreateMajorMeaningHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.MeaningMajorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateMeaningMajor(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateMajorMeaningHandler updates an existing MajorMeaning
// @Summary Update a MajorMeaning
// @Description Updates an existing MajorMeaning
// @Tags meanings
// @Accept json
// @Produce json
// @Param id path int true "MajorMeaning ID"
// @Param MajorMeaning body models.MeaningMajorInput true "Updated MajorMeaning"
// @Success 200 {object} models.MeaningMajor
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /meanings/major/{id} [put]
func UpdateMajorMeaningHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.MeaningMajorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		updated, err := models.UpdateMajorMeaning(a.DB, id, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Major Meaning not found")
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteMajorMeaningHandler deletes a MajorMeaning by ID
// @Summary Delete a MajorMeaning
// @Description Deletes a MajorMeaning by ID
// @Tags meanings
// @Param id path int true "MajorMeaning ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /meanings/major/{id} [delete]
func DeleteMajorMeaningHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteMajorMeaning(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
