package routes

import (
	"os"

	_ "github.com/ilbagatto/tarot-api/docs" // Import generated Swagger docs
	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/handlers"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// InitRoutes initializes all application routes
func InitRoutes(a *app.App) {
	e := a.Echo // Using Echo instance from App struct

	// Define API routes
	// Decks routes
	e.GET("/decks", handlers.ListDecksHandler(a))
	e.GET("/decks/:id", handlers.GetDeckByIDHandler(a))
	e.POST("/decks", handlers.CreateDeckHandler(a))
	e.PUT("/decks/:id", handlers.UpdateDeckHandler(a))
	e.DELETE("/decks/:id", handlers.DeleteDeckHandler(a))
	// Source routes
	e.GET("/sources", handlers.ListSourcesHandler(a))
	e.GET("/sources/:id", handlers.GetSourceByIDHandler(a))
	e.POST("/sources", handlers.CreateSourceHandler(a))
	e.PUT("/sources/:id", handlers.UpdateSourceHandler(a))
	e.DELETE("/sources/:id", handlers.DeleteSourceHandler(a))
	// Spreads
	e.GET("/spreads", handlers.ListSpreadsHandler(a))
	e.GET("/spreads/:id", handlers.GetSpreadByIDHandler(a))
	e.POST("/spreads", handlers.CreateSpreadHandler(a))
	e.PUT("/spreads/:id", handlers.UpdateSpreadHandler(a))
	e.DELETE("/spreads/:id", handlers.DeleteSpreadHandler(a))

	// Suits
	e.GET("/suits", handlers.ListSuitsHandler(a))
	e.GET("/suits/:id", handlers.GetSuitByIDHandler(a))
	e.POST("/suits", handlers.CreateSuitHandler(a))
	e.PUT("/suits/:id", handlers.UpdateSuitHandler(a))
	e.DELETE("/suits/:id", handlers.DeleteSuitHandler(a))

	// Ranks
	e.GET("/ranks", handlers.ListRanksHandler(a))
	e.GET("/ranks/:id", handlers.GetRankByIDHandler(a))
	e.POST("/ranks", handlers.CreateRankHandler(a))
	e.PUT("/ranks/:id", handlers.UpdateRankHandler(a))
	e.DELETE("/ranks/:id", handlers.DeleteRankHandler(a))

	// Major Arcana Cards
	e.GET("/cards/major", handlers.ListMajorCardsHandler(a))
	e.GET("/cards/major/:id", handlers.GetMajorCardByIDHandler(a))
	e.POST("/cards/major", handlers.CreateMajorCardHandler(a))
	e.PUT("/cards/major/:id", handlers.UpdateMajorCardHandler(a))
	e.DELETE("/cards/major/:id", handlers.DeleteMajorCardHandler(a))

	// Minor Arcana Cards
	e.GET("/cards/minor", handlers.ListMinorCardsHandler(a))
	e.GET("/cards/minor/:id", handlers.GetMinorCardByIDHandler(a))
	e.POST("/cards/minor", handlers.CreateMinorCardHandler(a))
	e.PUT("/cards/minor/:id", handlers.UpdateMinorCardHandler(a))
	e.DELETE("/cards/minor/:id", handlers.DeleteMinorCardHandler(a))

	// Major cards meanings
	e.GET("/meanings/major", handlers.ListMajorMeaningsHandler(a))
	e.GET("/meanings/major/:id", handlers.GetMajorMeaningByIDHandler(a))
	e.POST("/meanings/major", handlers.CreateMajorMeaningHandler(a))
	e.PUT("/meanings/major/:id", handlers.UpdateMajorMeaningHandler(a))
	e.DELETE("/meanings/major/:id", handlers.DeleteMajorMeaningHandler(a))

	// Minor cards meanings
	e.GET("/meanings/minor", handlers.ListMinorMeaningsHandler(a))
	e.GET("/meanings/minor/:id", handlers.GetMinorMeaningByIDHandler(a))
	e.POST("/meanings/minor", handlers.CreateMinorMeaningHandler(a))
	e.PUT("/meanings/minor/:id", handlers.UpdateMinorMeaningHandler(a))
	e.DELETE("/meanings/minor/:id", handlers.DeleteMinorMeaningHandler(a))

	// Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if os.Getenv("APP_ENV") == "dev" {
		e.GET("/*", func(c echo.Context) error {
			c.Response().Header().Del(echo.HeaderContentType)
			return c.File("static/images" + c.Request().URL.Path)
		})
	}

}
