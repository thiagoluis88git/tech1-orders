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

func TestGetOrderByIDHandler(t *testing.T) {
	t.Parallel()

	t.Run("got success when calling get order by id handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/orders/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getOrderByIdUseCase := new(MockGetOrderByIdUseCase)

		getOrderByIdUseCase.On("Execute", req.Context(), uint(12)).
			Return(dto.OrderResponse{
				OrderId: uint(12),
			}, nil)

		getOrderByIdHandler := handler.GetOrderByIdHandler(getOrderByIdUseCase)

		getOrderByIdHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.OrderResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NoError(t, err)

		assert.Equal(t, uint(12), response.OrderId)
	})

	t.Run("got error on GetOrderId UseCase when calling get order by id handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/orders/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getOrderByIdUseCase := new(MockGetOrderByIdUseCase)

		getOrderByIdUseCase.On("Execute", req.Context(), uint(12)).
			Return(dto.OrderResponse{}, &responses.BusinessResponse{
				StatusCode: 409,
			})

		getOrderByIdHandler := handler.GetOrderByIdHandler(getOrderByIdUseCase)

		getOrderByIdHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusConflict, recorder.Code)
	})

	t.Run("got error on invalid id param when calling get order by id handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/orders/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "x12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getOrderByIdUseCase := new(MockGetOrderByIdUseCase)

		getOrderByIdHandler := handler.GetOrderByIdHandler(getOrderByIdUseCase)

		getOrderByIdHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got error on missing id param when calling get order by id handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockOrder())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/orders/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getOrderByIdUseCase := new(MockGetOrderByIdUseCase)

		getOrderByIdHandler := handler.GetOrderByIdHandler(getOrderByIdUseCase)

		getOrderByIdHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
