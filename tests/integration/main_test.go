package integration

import (
	"os"
	"testing"

	"github.com/ilbagatto/tarot-api/internal/testutils"
)

var testApp *testutils.TestApp

func TestMain(m *testing.M) {
	testApp = testutils.SetupTestApp()
	code := m.Run()
	testApp.Close()
	os.Exit(code)
}
