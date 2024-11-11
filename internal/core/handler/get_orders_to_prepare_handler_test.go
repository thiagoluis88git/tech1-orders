package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/handler"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func TestGetOrdersWaitingPrepareHandler(t *testing.T) {
	t.Parallel()

	t.Run("got success when calling get orders to prepare handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/orders-to-prepare", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getOrdersToPrepareUseCase := new(MockGetOrdersToPrepareUseCase)

		getOrdersToPrepareUseCase.On("Execute", req.Context()).
			Return([]dto.OrderResponse{
				{
					OrderId: uint(12),
				},
			}, nil)

		getOrdersToPrepareHandler := handler.GetOrdersToPrepareHandler(getOrdersToPrepareUseCase)

		getOrdersToPrepareHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response []dto.OrderResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NoError(t, err)

		assert.Equal(t, uint(12), response[0].OrderId)
	})

	t.Run("got error on GetOrders Use Case when calling get orders to prepare handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/orders-to-prepare", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getOrdersToPrepareUseCase := new(MockGetOrdersToPrepareUseCase)

		getOrdersToPrepareUseCase.On("Execute", req.Context()).
			Return([]dto.OrderResponse{}, &responses.BusinessResponse{
				StatusCode: 422,
			})

		getOrdersToPrepareHandler := handler.GetOrdersToPrepareHandler(getOrdersToPrepareUseCase)

		getOrdersToPrepareHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})
}
