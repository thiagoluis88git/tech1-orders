package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func TestOrderValidations(t *testing.T) {
	t.Parallel()

	t.Run("got success when validating order to prepare with status CRIADO", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToPrepareUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderId:     uint(1),
			OrderStatus: "Criado",
		}, nil)

		err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
	})

	t.Run("got error when validating order to prepare with status different than CRIADO", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToPrepareUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderId:     uint(1),
			OrderStatus: "Finalizado",
		}, nil)

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)
	})

	t.Run("got error when validating order to prepare with repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToPrepareUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{}, &responses.NetworkError{
			Code: 409,
		})

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)
	})

	t.Run("got success when validating order to done with status PREPARANDO", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToDoneUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderId:     uint(1),
			OrderStatus: "Preparando",
		}, nil)

		err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
	})

	t.Run("got error when validating order to done with status different than PREPARANDO", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToDoneUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderId:     uint(1),
			OrderStatus: "Finalizado",
		}, nil)

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)
	})

	t.Run("got error when validating order to done with repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToDoneUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{}, &responses.NetworkError{
			Code: 409,
		})

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)
	})

	t.Run("got success when validating order to delivered or not with status FINALIZADO", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderId:     uint(1),
			OrderStatus: "Finalizado",
		}, nil)

		err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
	})

	t.Run("got error when validating order to delivered or not with status different than FINALIZADO", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderId:     uint(1),
			OrderStatus: "Preparando",
		}, nil)

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)
	})

	t.Run("got error when validating order to delivered or not with repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{}, &responses.NetworkError{
			Code: 409,
		})

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)
	})
}
