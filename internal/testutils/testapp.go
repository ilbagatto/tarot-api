package testutils

import (
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"

	"github.com/ilbagatto/tarot-api/internal/app"
	"github.com/ilbagatto/tarot-api/internal/db"
	"github.com/ilbagatto/tarot-api/internal/routes"
	"github.com/joho/godotenv"
)

type TestApp struct {
	App *app.App
}

func loadEnvFromProjectRoot() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get current file location")
	}

	projectRoot := filepath.Join(filepath.Dir(filename), "../../")
	envPath := filepath.Join(projectRoot, ".env.test")

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️ Could not load %s: %v", envPath, err)
	}
}

// SetupTestApp initializes the application for integration tests
func SetupTestApp() *TestApp {
	loadEnvFromProjectRoot()
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}
	a := app.NewApp(database)
	routes.InitRoutes(a)
	return &TestApp{App: a}
}

// Close shuts down the test app and closes the database
func (ta *TestApp) Close() {
	if ta.App.DB != nil {
		_ = ta.App.DB.Close()
	}
}

// Request makes a test HTTP request using Echo and returns the response
func (ta *TestApp) Request(method, path string, body http.Handler) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	ta.App.Echo.ServeHTTP(rec, req)
	return rec
}
