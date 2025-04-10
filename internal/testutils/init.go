package testutils

import (
	"database/sql"
	"log"
	"testing"

	"github.com/ilbagatto/tarot-api/internal/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func LoadTestEnv() {
	if err := godotenv.Load(".env.test"); err != nil {
		log.Println("⚠️  Could not load .env.test:", err)
	}
}

// InitTestDBWithAssert loads env and returns connected DB with assertion
func InitTestDBWithAssert(t *testing.T) *sql.DB {
	err := godotenv.Load(".env.test")
	assert.NoError(t, err, "failed to load .env.test")

	database, err := db.InitDB()
	assert.NoError(t, err, "should connect to test DB without error")

	return database
}
