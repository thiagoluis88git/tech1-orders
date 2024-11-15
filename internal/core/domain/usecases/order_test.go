package usecases

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func mockCustomer() dto.Customer {
	return dto.Customer{
		ID:   uint(1),
		Name: "CustomerName",
	}
}

func TestOrderServices(t *testing.T) {
	t.Parallel()

	t.Run("got success when generating ticket number in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		customerRepo := new(MockCustomerRepository)
		validateToPrepare := NewValidateOrderToPrepareUseCase(mockRepo)
		validateToDone := NewValidateOrderToDoneUseCase(mockRepo)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewCreateOrderUseCase(
			mockRepo,
			customerRepo,
			validateToPrepare,
			validateToDone,
			validateToDeliveredOrNot,
			sortOrdersUseCase,
		)

		ctx := context.TODO()

		date := time.Now().UnixMilli()

		mockRepo.On("GetNextTicketNumber", ctx, date).Return(1, nil)

		response := sut.GenerateTicket(ctx, date)

		assert.Equal(t, 1, response)
	})

	t.Run("got success when creating order in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		customerRepo := new(MockCustomerRepository)
		validateToPrepare := NewValidateOrderToPrepareUseCase(mockRepo)
		validateToDone := NewValidateOrderToDoneUseCase(mockRepo)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewCreateOrderUseCase(
			mockRepo,
			customerRepo,
			validateToPrepare,
			validateToDone,
			validateToDeliveredOrNot,
			sortOrdersUseCase,
		)

		ctx := context.TODO()

		date := time.Now().UnixMilli()

		mockRepo.On("CreateOrder", ctx, orderCreation).Return(orderCreationResponse, nil)
		mockRepo.On("GetNextTicketNumber", ctx, date).Return(1, nil)

		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)

		wg.Add(1)
		response, err := sut.Execute(ctx, orderCreation, date, wg, ch)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)
	})

	t.Run("got success when creating order with customer in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		customerRepo := new(MockCustomerRepository)
		validateToPrepare := NewValidateOrderToPrepareUseCase(mockRepo)
		validateToDone := NewValidateOrderToDoneUseCase(mockRepo)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewCreateOrderUseCase(
			mockRepo,
			customerRepo,
			validateToPrepare,
			validateToDone,
			validateToDeliveredOrNot,
			sortOrdersUseCase,
		)

		ctx := context.TODO()

		date := time.Now().UnixMilli()

		customerRepo.On("GetCustomerByCPF", ctx, *orderCreationWithCustomer.CPF).Return(mockCustomer(), nil)
		mockRepo.On("CreateOrder", ctx, orderCreationWithCustomer).Return(orderWithCustomerCreationResponse, nil)
		mockRepo.On("GetNextTicketNumber", ctx, date).Return(1, nil)

		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)

		wg.Add(1)
		response, err := sut.Execute(ctx, orderCreationWithCustomer, date, wg, ch)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)
	})

	t.Run("got error when creating order in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		customerRepo := new(MockCustomerRepository)
		validateToPrepare := NewValidateOrderToPrepareUseCase(mockRepo)
		validateToDone := NewValidateOrderToDoneUseCase(mockRepo)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewCreateOrderUseCase(
			mockRepo,
			customerRepo,
			validateToPrepare,
			validateToDone,
			validateToDeliveredOrNot,
			sortOrdersUseCase,
		)

		ctx := context.TODO()

		date := time.Now().UnixMilli()

		mockRepo.On("GetNextTicketNumber", ctx, date).Return(1, nil)
		mockRepo.On("CreateOrder", ctx, orderCreationWithCustomer).Return(dto.OrderResponse{}, &responses.NetworkError{
			Code:    409,
			Message: "Conflict",
		})

		wg := &sync.WaitGroup{}
		ch := make(chan bool, 1)

		wg.Add(1)
		response, err := sut.Execute(ctx, orderCreationWithCustomer, date, wg, ch)

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusConflict, businessError.StatusCode)
	})

	t.Run("got success when getting order by id in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)

		sut := NewGetOrderByIdUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(orderCreationResponse, nil)

		response, err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
		assert.NotEmpty(t, response)
	})

	t.Run("got error when getting order by id in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)

		sut := NewGetOrderByIdUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderProductResponse{}, &responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		response, err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when getting orders to prepare in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewGetOrdersToPrepareUseCase(mockRepo, sortOrdersUseCase)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Criado",
		}, nil)
		mockRepo.On("GetOrdersToPrepare", ctx).Return(ordersList, nil)

		response, err := sut.Execute(ctx)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)

		assert.Equal(t, 1, len(response))
	})

	t.Run("got error when getting orders to prepare in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewGetOrdersToPrepareUseCase(mockRepo, sortOrdersUseCase)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Criado",
		}, nil)
		mockRepo.On("GetOrdersToPrepare", ctx).Return(dto.OrderResponse{}, &responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		response, err := sut.Execute(ctx)

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when getting orders status in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewGetOrdersToFollowUseCase(mockRepo, sortOrdersUseCase)

		ctx := context.TODO()

		mockRepo.On("GetOrdersToFollow", ctx).Return(ordersList, nil)

		response, err := sut.Execute(ctx)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)

		assert.Equal(t, 1, len(response))
	})

	t.Run("got error when getting orders status in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sortOrdersUseCase := NewSortOrdersUseCase()

		sut := NewGetOrdersToFollowUseCase(mockRepo, sortOrdersUseCase)

		ctx := context.TODO()

		mockRepo.On("GetOrdersToFollow", ctx).Return(dto.OrderResponse{}, &responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		response, err := sut.Execute(ctx)

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when updating order to delivered in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)

		sut := NewUpdateToDeliveredUseCase(mockRepo, validateToDeliveredOrNot)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Finalizado",
		}, nil)
		mockRepo.On("UpdateToDelivered", ctx, uint(1)).Return(nil)

		err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
	})

	t.Run("got error when updating order to delivered in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)

		sut := NewUpdateToDeliveredUseCase(mockRepo, validateToDeliveredOrNot)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Finalizado",
		}, nil)
		mockRepo.On("UpdateToDelivered", ctx, uint(1)).Return(&responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when updating order to done in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToDone := NewValidateOrderToDoneUseCase(mockRepo)

		sut := NewUpdateToDoneUseCase(mockRepo, validateToDone)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Preparando",
		}, nil)
		mockRepo.On("UpdateToDone", ctx, uint(1)).Return(nil)

		err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
	})

	t.Run("got error when updating order to done in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToDone := NewValidateOrderToDoneUseCase(mockRepo)

		sut := NewUpdateToDoneUseCase(mockRepo, validateToDone)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Preparando",
		}, nil)
		mockRepo.On("UpdateToDone", ctx, uint(1)).Return(&responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when updating order to not delivered in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)

		sut := NewUpdateToNotDeliveredUseCase(mockRepo, validateToDeliveredOrNot)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Finalizado",
		}, nil)
		mockRepo.On("UpdateToNotDelivered", ctx, uint(1)).Return(nil)

		err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
	})

	t.Run("got error when updating order to not delivered in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToDeliveredOrNot := NewValidateOrderToDeliveredOrNotUseCase(mockRepo)

		sut := NewUpdateToNotDeliveredUseCase(mockRepo, validateToDeliveredOrNot)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Finalizado",
		}, nil)
		mockRepo.On("UpdateToNotDelivered", ctx, uint(1)).Return(&responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when updating order to preparing in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToPrepare := NewValidateOrderToPrepareUseCase(mockRepo)

		sut := NewUpdateToPreparingUseCase(mockRepo, validateToPrepare)

		ctx := context.TODO()

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Criado",
		}, nil)
		mockRepo.On("UpdateToPreparing", ctx, uint(1)).Return(nil)

		err := sut.Execute(ctx, uint(1))

		assert.NoError(t, err)
	})

	t.Run("got error when updating order to preparing in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		validateToPrepare := NewValidateOrderToPrepareUseCase(mockRepo)

		sut := NewUpdateToPreparingUseCase(mockRepo, validateToPrepare)

		ctx := context.TODO()

		mockRepo.On("UpdateToPreparing", ctx, uint(1)).Return(&responses.NetworkError{
			Code:    404,
			Message: "Not Found",
		})

		mockRepo.On("GetOrderById", ctx, uint(1)).Return(dto.OrderResponse{
			OrderStatus: "Criado",
		}, nil)
		err := sut.Execute(ctx, uint(1))

		assert.Error(t, err)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusNotFound, businessError.StatusCode)
	})

	t.Run("got success when get orders waiting payment use case", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewGetOrdersWaitingPaymentUseCase(mockRepo, NewSortOrdersUseCase())

		ctx := context.TODO()

		mockRepo.On("GetOrdersWaitingPayment", ctx).Return([]dto.OrderResponse{
			{
				OrderId:      uint(2),
				TicketNumber: 12,
			},
			{
				OrderId:      uint(3),
				TicketNumber: 13,
			},
		}, nil)

		response, err := sut.Execute(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
	})

	t.Run("got error on GetOrdersWaitingPayment Repository when get orders waiting payment use case", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockOrderRepository)
		sut := NewGetOrdersWaitingPaymentUseCase(mockRepo, NewSortOrdersUseCase())

		ctx := context.TODO()

		mockRepo.On("GetOrdersWaitingPayment", ctx).Return([]dto.OrderResponse{}, &responses.LocalError{
			Code: 409,
		})

		response, err := sut.Execute(ctx)

		assert.Error(t, err)
		assert.Empty(t, response)
	})
}
