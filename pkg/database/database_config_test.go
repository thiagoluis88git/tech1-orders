package database_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/pkg/database"
	"github.com/thiagoluis88git/tech1-orders/pkg/environment"
	"gorm.io/driver/postgres"
)

func setup() {
	os.Setenv(environment.QRCodeGatewayRootURL, "ROOT_URL")
	os.Setenv(environment.DBHost, "HOST")
	os.Setenv(environment.DBPort, "1234")
	os.Setenv(environment.DBUser, "User")
	os.Setenv(environment.DBPassword, "Pass")
	os.Setenv(environment.DBName, "Name")
	os.Setenv(environment.WebhookMercadoLivrePaymentURL, "WEBHOOK")
	os.Setenv(environment.QRCodeGatewayToken, "token")
	os.Setenv(environment.Region, "Region")
	os.Setenv(environment.CustomerRootAPI, "CustomerRootAPI")
}

func TestDatabaseConfig(t *testing.T) {
	t.Parallel()
	setup()

	t.Run("got success when starting config database", func(t *testing.T) {
		environment.LoadEnvironmentVariables()

		conn, _, err := sqlmock.New()
		assert.NoError(t, err)

		dialector := postgres.New(postgres.Config{
			DSN:                  "sqlmock_db_0",
			DriverName:           "postgres",
			Conn:                 conn,
			PreferSimpleProtocol: true,
		})

		config, err := database.ConfigDatabase(dialector)

		assert.NoError(t, err)
		assert.NotEmpty(t, config)
	})

	t.Run("got error when starting config database", func(t *testing.T) {
		environment.LoadEnvironmentVariables()

		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v",
			environment.GetDBHost(),
			environment.GetDBUser(),
			environment.GetDBPassword(),
			environment.GetDBName(),
			environment.GetDBPort(),
		)

		config, err := database.ConfigDatabase(postgres.Open(dsn))

		assert.Error(t, err)
		assert.Empty(t, config)
	})
}
