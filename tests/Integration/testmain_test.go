package integration

import (
	"os"
	"testing"

	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
)

func TestMain(m *testing.M) {
    setup.InitTestDatabase()
    defer setup.TeardownTestDB()
    code := m.Run()
    os.Exit(code)
}
