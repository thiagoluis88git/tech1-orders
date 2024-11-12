package httpserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/pkg/httpserver"
)

func TestHttpClient(t *testing.T) {
	t.Parallel()

	t.Run("got success when creating http client", func(t *testing.T) {
		t.Parallel()

		client := httpserver.NewHTTPClient()

		assert.NotEmpty(t, client)
	})

	t.Run("got success when calling DoRequest client", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpserver.SendResponseSuccess(w, dto.ProductResponse{
				Id:   uint(1),
				Name: "Name",
			})
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		response, err := httpserver.DoGetRequest(context.TODO(), ts.Client(), ts.URL, nil, dto.ProductResponse{})

		assert.NoError(t, err)
		assert.NotEmpty(t, response)
	})

	t.Run("got error on invalid json when calling DoRequest client", func(t *testing.T) {
		t.Parallel()

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		ts := httptest.NewServer(responseHandler)
		defer ts.Close()

		response, err := httpserver.DoGetRequest(context.TODO(), ts.Client(), ts.URL, nil, dto.ProductResponse{})

		assert.Error(t, err)
		assert.Empty(t, response)
	})

	t.Run("got error when calling DoRequest client", func(t *testing.T) {
		t.Parallel()

		client := httpserver.NewHTTPClient()

		assert.NotEmpty(t, client)

		response, err := httpserver.DoGetRequest(context.TODO(), client, "http://localhost", nil, dto.ProductResponse{})

		assert.Error(t, err)
		assert.Empty(t, response)
	})
}
