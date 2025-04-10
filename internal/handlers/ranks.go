package handlers

import (
	"net/http"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/models"
	"github.com/labstack/echo/v4"
)

// ListRanksHandler returns all ranks
// @Summary Get all ranks
// @Description Retrieves a list of all available interpretation ranks
// @Tags ranks
// @Produce json
// @Success 200 {array} models.Rank
// @Failure 500 {object} APIResponse
// @Router /ranks [get]
func ListRanksHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		ranks, err := models.ListRanks(a.DB)
		if err != nil {
			return useHandleDBError(c, err)
		}
		return c.JSON(http.StatusOK, ranks)
	}
}

// GetRankByIDHandler returns a specific rank by ID
// @Summary Get rank by ID
// @Description Retrieves a rank by its ID
// @Tags ranks
// @Produce json
// @Param id path int true "Rank ID"
// @Success 200 {object} models.Rank
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Router /ranks/{id} [get]
func GetRankByIDHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		src, err := models.GetRankByID(a.DB, id)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Rank not found")
		}

		return c.JSON(http.StatusOK, src)
	}
}

// CreateRankHandler creates a new rank
// @Summary Create a new rank
// @Description Adds a new rank of interpretations
// @Tags ranks
// @Accept json
// @Produce json
// @Param rank body models.RankInput true "Rank data"
// @Success 201 {object} models.IDOnly "Returns the ID of the created rank"
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /ranks [post]
func CreateRankHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.RankInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		id, err := models.CreateRank(a.DB, input)
		if err != nil {
			return useHandleDBError(c, err)
		}

		return c.JSON(http.StatusCreated, models.IDOnly{ID: *id})
	}
}

// UpdateRankHandler updates an existing rank
// @Summary Update a rank
// @Description Updates an existing rank
// @Tags ranks
// @Accept json
// @Produce json
// @Param id path int true "Rank ID"
// @Param rank body models.RankInput true "Updated rank"
// @Success 200 {object} models.Rank
// @Failure 400 {object} handlers.APIResponse
// @Failure 404 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /ranks/{id} [put]
func UpdateRankHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam(c)
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		var input models.RankInput
		if err := useBind(c, &input); err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		updated, err := models.UpdateRank(a.DB, id, input)
		if err != nil {
			return useHandleNotFoundOrDBError(c, err, "Rank not found")
		}

		return c.JSON(http.StatusOK, updated)
	}
}

// DeleteRankHandler deletes a rank by ID
// @Summary Delete a rank
// @Description Deletes a rank by ID
// @Tags ranks
// @Param id path int true "Rank ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /ranks/{id} [delete]
func DeleteRankHandler(a *app.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := useIDParam((c))
		if err != nil {
			return SendError(c, http.StatusBadRequest, err)
		}
		if err := models.DeleteRank(a.DB, id); err != nil {
			return useHandleDBError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
