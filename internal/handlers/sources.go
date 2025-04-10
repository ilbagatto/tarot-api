package handlers

import (
	"net/http"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListSourcesHandler returns all sources
// @Summary Get all sources
// @Description Retrieves a list of all available interpretation sources
// @Tags sources
// @Produce json
// @Success 200 {array} models.SourceListItem
// @Failure 500 {object} APIResponse
// @Router /sources [get]
func ListSourcesHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		sources, err := models.ListSources(a.DB)
		if err != nil {
			return useHandleDBError(c, err)
		}
		return c.JSON(http.StatusOK, sources)
	}
}

// GetSourceByIDHandler returns a specific source by ID
// @Summary Get source by ID
// @Description Retrieves a source by its ID
// @Tags sources
// @Produce json
// @Param id path int true "Source ID"
// @Success 200 {object} models.Source
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /sources/{id} [get]
func GetSourceByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetSourceByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Source not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateSourceHandler creates a new source
// @Summary Create a new source
// @Description Adds a new source of interpretations
// @Tags sources
// @Accept json
// @Produce json
// @Param source body models.SourceInput true "Source data with list of deck IDs"
// @Success 201 {object} models.IDOnly "Returns the ID of the created source"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /sources [post]
func CreateSourceHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.SourceInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateSource(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateSourceHandler updates an existing source
// @Summary Update a source
// @Description Updates an existing source
// @Tags sources
// @Accept json
// @Produce json
// @Param id path int true "Source ID"
// @Param source body models.SourceInput true "Updated source with list of deck IDs"
// @Success 200 {object} models.Source
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /sources/{id} [put]
func UpdateSourceHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.SourceInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		updated, err := models.UpdateSource(a.DB, id, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Source not found")
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteSourceHandler deletes a source by ID
// @Summary Delete a source
// @Description Deletes a source by ID
// @Tags sources
// @Param id path int true "Source ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /sources/{id} [delete]
func DeleteSourceHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteSource(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
