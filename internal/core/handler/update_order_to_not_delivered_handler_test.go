package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/handler"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func TestUpdateOrderToNotDeliveredHandler(t *testing.T) {
	t.Parallel()

	t.Run("got success when calling update to not delivered handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPut, "/api/orders/{id}/not-delivered", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateOrderToNotDelivered := new(MockUpdateToNotDeliveredUseCase)

		updateOrderToNotDelivered.On("Execute", req.Context(), uint(12)).
			Return(dto.OrderResponse{
				OrderId: uint(12),
			}, nil)

		updateOrderToNotDeliveredHandler := handler.UpdateOrderNotDeliveredandler(updateOrderToNotDelivered)

		updateOrderToNotDeliveredHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("got error on UpdateToNotDelivered UseCase when calling update to not delivered handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPut, "/api/orders/{id}/not-delivered", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateOrderToNotDelivered := new(MockUpdateToNotDeliveredUseCase)

		updateOrderToNotDelivered.On("Execute", req.Context(), uint(12)).
			Return(dto.OrderResponse{}, &responses.BusinessResponse{
				StatusCode: 422,
			})

		updateOrderToNotDeliveredHandler := handler.UpdateOrderNotDeliveredandler(updateOrderToNotDelivered)

		updateOrderToNotDeliveredHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})

	t.Run("got error on invalid id when calling update to not delivered handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPut, "/api/orders/{id}/not-delivered", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "x12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateOrderToNotDelivered := new(MockUpdateToNotDeliveredUseCase)

		updateOrderToNotDelivered.On("Execute", req.Context(), uint(12)).
			Return(dto.OrderResponse{}, &responses.BusinessResponse{
				StatusCode: 422,
			})

		updateOrderToNotDeliveredHandler := handler.UpdateOrderNotDeliveredandler(updateOrderToNotDelivered)

		updateOrderToNotDeliveredHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got error on missing id when calling update to not delivered handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPut, "/api/orders/{id}/not-delivered", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateOrderToNotDelivered := new(MockUpdateToNotDeliveredUseCase)

		updateOrderToNotDelivered.On("Execute", req.Context(), uint(12)).
			Return(dto.OrderResponse{}, &responses.BusinessResponse{
				StatusCode: 422,
			})

		updateOrderToNotDeliveredHandler := handler.UpdateOrderNotDeliveredandler(updateOrderToNotDelivered)

		updateOrderToNotDeliveredHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
