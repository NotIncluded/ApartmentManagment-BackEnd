package integration

import (
	"os"
	"testing"

	"github.com/PunMung-66/ApartmentSys/tests/Integration/setup"
)

func TestMain(m *testing.M) {
	setup.InitTestDatabase(setup.Env)
	defer setup.TeardownTestDB()
	code := m.Run()
	os.Exit(code)
}
