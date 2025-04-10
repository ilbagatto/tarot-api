package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListMinorMeaningsHandler returns filtered interpretations for Minor Arcana cards
// @Summary Get Minor Arcana meanings
// @Description Returns a list of meanings for minor arcana cards with optional filters
// @Tags meanings
// @Accept json
// @Produce json
// @Param number query int false "Card number (optional)"
// @Param position query string false "Card position" Enums(straight, reverted)
// @Param source query int false "Source ID (optional)"
// @Success 200 {array} models.MeaningMinor
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /meanings/minor [get]
func ListMinorMeaningsHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		filters := map[string]any{}

		if source := c.QueryParam("suit"); source != "" {
			if s, err := strconv.ParseInt(source, 10, 64); err == nil {
				filters["suit"] = s
			} else {
				return SendError(c, http.StatusBadRequest, fmt.Errorf("invalid suit"))
			}
		}
		if source := c.QueryParam("rank"); source != "" {
			if s, err := strconv.ParseInt(source, 10, 64); err == nil {
				filters["rank"] = s
			} else {
				return SendError(c, http.StatusBadRequest, fmt.Errorf("invalid rank"))
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

		result, err := models.ListMinorMeanings(a.DB, filters)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusOK, result)
	}
}

// GetMinorMeaningByIDHandler returns a specific MinorMeaning by ID
// @Summary Get MinorMeaning by ID
// @Description Retrieves a MinorMeaning by its ID
// @Tags meanings
// @Produce json
// @Param id path int true "MinorMeaning ID"
// @Success 200 {object} models.MeaningMinor
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /meanings/minor/{id} [get]
func GetMinorMeaningByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetMinorMeaningByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Minor Meaning not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateMinorMeaningHandler creates a new MinorMeaning
// @Summary Create a new MinorMeaning
// @Description Adds a new MinorMeaning of interpretations
// @Tags meanings
// @Accept json
// @Produce json
// @Param MinorMeaning body models.MeaningMinorInput true "MinorMeaning data"
// @Success 201 {object} models.IDOnly "Returns the ID of the created MinorMeaning"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /meanings/minor [post]
func CreateMinorMeaningHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.MeaningMinorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateMinorMeaning(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateMinorMeaningHandler updates an existing MinorMeaning
// @Summary Update a MinorMeaning
// @Description Updates an existing MinorMeaning
// @Tags meanings
// @Accept json
// @Produce json
// @Param id path int true "MinorMeaning ID"
// @Param MinorMeaning body models.MeaningMinorInput true "Updated MinorMeaning"
// @Success 200 {object} models.MeaningMinor
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /meanings/minor/{id} [put]
func UpdateMinorMeaningHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.MeaningMinorInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		updated, err := models.UpdateMinorMeaning(a.DB, id, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Minor Meaning not found")
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteMinorMeaningHandler deletes a MinorMeaning by ID
// @Summary Delete a MinorMeaning
// @Description Deletes a MinorMeaning by ID
// @Tags meanings
// @Param id path int true "MinorMeaning ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /meanings/minor/{id} [delete]
func DeleteMinorMeaningHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteMinorMeaning(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
