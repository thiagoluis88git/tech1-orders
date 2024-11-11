package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/handler"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func mockOrder() dto.Order {
	return dto.Order{
		TotalPrice: 123.34,
		PaymentID:  uint(4),
		OrderProduct: []dto.OrderProduct{
			{
				ProductID:    1,
				ProductPrice: 23.4,
			},
			{
				ProductID:    2,
				ProductPrice: 23.4,
			},
		},
	}
}

func TestCreateOrderHandler(t *testing.T) {
	t.Parallel()

	t.Run("got success when calling create order handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/order", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		createOrderUseCase := new(MockCreateOrderUseCase)

		now := time.Now()
		orderDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		var wg sync.WaitGroup
		channel := make(chan bool, 1)
		wg.Add(1)

		createOrderUseCase.On("Execute", req.Context(), mockOrder(), orderDate.UnixMilli(), &wg, channel).
			Return(dto.OrderResponse{
				OrderId: uint(2),
			}, nil)

		loginCustomerHandler := handler.CreateOrderHandler(createOrderUseCase)

		loginCustomerHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.OrderResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NoError(t, err)

		assert.Equal(t, uint(2), response.OrderId)
	})

	t.Run("got error on CreateOrder UseCase when calling create order handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/order", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		createOrderUseCase := new(MockCreateOrderUseCase)

		now := time.Now()
		orderDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		var wg sync.WaitGroup
		channel := make(chan bool, 1)
		wg.Add(1)

		createOrderUseCase.On("Execute", req.Context(), mockOrder(), orderDate.UnixMilli(), &wg, channel).
			Return(dto.OrderResponse{}, &responses.BusinessResponse{
				StatusCode: 500,
			})

		loginCustomerHandler := handler.CreateOrderHandler(createOrderUseCase)

		loginCustomerHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("got error on Decoding invalid json when calling create order handler", func(t *testing.T) {
		t.Parallel()

		body := bytes.NewBuffer([]byte("sfgg{[]}"))

		req := httptest.NewRequest(http.MethodPost, "/api/order", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		createOrderUseCase := new(MockCreateOrderUseCase)

		loginCustomerHandler := handler.CreateOrderHandler(createOrderUseCase)

		loginCustomerHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
