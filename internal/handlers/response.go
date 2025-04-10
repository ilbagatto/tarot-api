package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// APIResponse defines a standard JSON response format
type APIResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ErrorResponse defines the standard error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"Duplicate entry"`
}

// NewSuccessResponse creates a success response with a message
func NewSuccessResponse(message string) APIResponse {
	return APIResponse{Message: message}
}

// NewErrorResponse creates an error response with an error message
func NewErrorResponse(err error) APIResponse {
	return APIResponse{Error: err.Error()}
}

// SendError sends a JSON error response with a status code
func SendError(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, NewErrorResponse(err))

}

// SendSuccess sends a JSON success response with a status code
func SendSuccess(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, NewSuccessResponse(message))
}

// HTTPErrorFromDBError analyzes a DB error and returns HTTP status code and APIResponse
func HTTPErrorFromDBError(err error) (int, APIResponse) {
	if err == nil {
		return http.StatusOK, APIResponse{}
	}

	msg := err.Error()

	switch {
	case strings.Contains(msg, "duplicate key"):
		return http.StatusConflict, APIResponse{Error: "Duplicate entry: " + msg}
	case strings.Contains(msg, "violates foreign key constraint"):
		return http.StatusConflict, APIResponse{Error: "Invalid reference: " + msg}
	case strings.Contains(msg, "invalid input syntax"):
		return http.StatusBadRequest, APIResponse{Error: "Invalid input: " + msg}
	default:
		return http.StatusInternalServerError, APIResponse{Error: msg}
	}
}
