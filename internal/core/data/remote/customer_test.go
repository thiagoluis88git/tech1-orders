package remote_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/remote"
	"github.com/thiagoluis88git/tech1-orders/pkg/environment"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
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
	os.Setenv(environment.CustomerRootAPI, "rootURL")
}

func TestCustomerRemote(t *testing.T) {
	t.Parallel()
	setup()

	t.Run("got success when getting customer by cpf remote", func(t *testing.T) {
		t.Parallel()

		environment.LoadEnvironmentVariables()

		recorder := httptest.NewRecorder()
		recorder.Header().Add("Content-Type", "application/json")
		recorder.WriteHeader(http.StatusOK)
		_, err := recorder.WriteString(MockCustomer)

		assert.NoError(t, err)

		resultExpected := recorder.Result()

		mockClient := &http.Client{
			Transport: &MockRoundTripper{
				Response: resultExpected,
			},
		}

		ds := remote.NewCustomerRemoteDataSource(mockClient, "rootURL")

		response, err := ds.GetCustomerByCPF(context.TODO(), "12345678910")

		assert.NoError(t, err)
		assert.Equal(t, "Customer Name", response.Name)
	})

	t.Run("got success when getting customer by cpf remote", func(t *testing.T) {
		t.Parallel()

		environment.LoadEnvironmentVariables()

		recorder := httptest.NewRecorder()
		recorder.Header().Add("Content-Type", "application/json")
		recorder.WriteHeader(http.StatusOK)
		_, err := recorder.WriteString("sdd{{}")

		assert.NoError(t, err)

		resultExpected := recorder.Result()

		mockClient := &http.Client{
			Transport: &MockRoundTripper{
				Response: resultExpected,
			},
		}

		ds := remote.NewCustomerRemoteDataSource(mockClient, "rootURL")

		response, err := ds.GetCustomerByCPF(context.TODO(), "12345678910")

		assert.Error(t, err)
		assert.Empty(t, response.Name)

		var netError *responses.NetworkError
		isNetError := errors.As(err, &netError)
		assert.Equal(t, true, isNetError)
		assert.Equal(t, http.StatusUnprocessableEntity, netError.Code)
	})
}
