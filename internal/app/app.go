package app

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

// App holds dependencies for the application
type App struct {
	DB   *sql.DB
	Echo *echo.Echo
}

func NewApp(db *sql.DB) *App {
	return &App{
		Echo: echo.New(),
		DB:   db,
	}
}
