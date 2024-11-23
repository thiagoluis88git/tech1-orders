package environment_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/pkg/environment"
)

func setup() {
	os.Setenv(environment.DBHost, "DBHost")
	os.Setenv(environment.DBPassword, "DBPassword")
	os.Setenv(environment.DBName, "DBName")
	os.Setenv(environment.DBPort, "DBPort")
	os.Setenv(environment.DBUser, "DBUser")
	os.Setenv(environment.Region, "Region")
	os.Setenv(environment.CustomerRootAPI, "CustomerRootAPI")
}

func TestEnvironment(t *testing.T) {
	t.Parallel()
	setup()

	t.Run("got success when loading variables", func(t *testing.T) {
		environment.LoadEnvironmentVariables()
	})

	t.Run("got success when initializing environment", func(t *testing.T) {
		environment.LoadEnvironmentVariables()

		assert.Equal(t, "DBHost", environment.GetDBHost())
		assert.Equal(t, "DBPassword", environment.GetDBPassword())
		assert.Equal(t, "DBPort", environment.GetDBPort())
		assert.Equal(t, "DBName", environment.GetDBName())
		assert.Equal(t, "DBUser", environment.GetDBUser())
		assert.Equal(t, "Region", environment.GetRegion())
		assert.Equal(t, "CustomerRootAPI", environment.GetCustomerRootAPI())
	})
}
