package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// useIDParam extracts and validates an integer ID from path or query parameters.
// If no parameter name is provided, defaults to "id".
// The parameter name is case-insensitive.
func useIDParam(c echo.Context, paramName ...string) (int64, error) {
	name := "id"
	if len(paramName) > 0 && paramName[0] != "" {
		name = paramName[0]
	}

	// 1. сначала ищем в path
	for i, key := range c.ParamNames() {
		if strings.EqualFold(key, name) {
			return strconv.ParseInt(c.ParamValues()[i], 10, 64)
		}
	}

	// 2. затем — в query
	val := c.QueryParam(name)
	if val == "" {
		return 0, fmt.Errorf("missing parameter: %s", name)
	}

	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be an integer", name)
	}
	return id, nil
}

// useBind binds and validates the request body into a target struct
func useBind[T any](c echo.Context, target *T) error {
	if err := c.Bind(target); err != nil {
		return errors.New("invalid request body")
	}
	return nil
}

// useHandleDBError translates model-level errors into HTTP responses
func useHandleDBError(c echo.Context, err error) error {
	if err == nil {
		return nil
	}
	status, resp := HTTPErrorFromDBError(err)
	return c.JSON(status, resp)
}

// useHandleNotFoundOrDBError handles sql.ErrNoRows and general DB errors
func useHandleNotFoundOrDBError(c echo.Context, err error, notFoundMsg string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return c.JSON(http.StatusNotFound, APIResponse{Error: notFoundMsg})
	}
	return useHandleDBError(c, err)
}
