package handler_test

import (
	"context"
	"sync"

	"github.com/stretchr/testify/mock"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

type MockPayOrderUseCase struct {
	mock.Mock
}

type MockGetPaymentTypesUseCase struct {
	mock.Mock
}

type MockCreateOrderUseCase struct {
	mock.Mock
}

type MockGetOrderByIdUseCase struct {
	mock.Mock
}

type MockGetOrdersToPrepareUseCase struct {
	mock.Mock
}

type MockGetOrdersToFollowUseCase struct {
	mock.Mock
}

type MockGetOrdersWaitingPaymentUseCase struct {
	mock.Mock
}

type MockUpdateToPreparingUseCase struct {
	mock.Mock
}

type MockUpdateToDoneUseCase struct {
	mock.Mock
}

type MockUpdateToDeliveredUseCase struct {
	mock.Mock
}

type MockUpdateToNotDeliveredUseCase struct {
	mock.Mock
}

type MockCreateProductUseCase struct {
	mock.Mock
}

type MockGetProductsByCategoryUseCase struct {
	mock.Mock
}

type MockGetProductsByIDUseCase struct {
	mock.Mock
}

type MockDeleteProductUseCase struct {
	mock.Mock
}

type MockUpdateProductUseCase struct {
	mock.Mock
}

type MockGetCategoryUseCase struct {
	mock.Mock
}

func (mock *MockPayOrderUseCase) Execute(ctx context.Context, payment dto.Payment) (dto.PaymentResponse, error) {
	args := mock.Called(ctx, payment)
	err := args.Error(1)

	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return args.Get(0).(dto.PaymentResponse), nil
}

func (mock *MockGetOrdersToPrepareUseCase) Execute(ctx context.Context) ([]dto.OrderResponse, error) {
	args := mock.Called(ctx)
	err := args.Error(1)

	if err != nil {
		return []dto.OrderResponse{}, err
	}

	return args.Get(0).([]dto.OrderResponse), nil
}

func (mock *MockGetOrdersToFollowUseCase) Execute(ctx context.Context) ([]dto.OrderResponse, error) {
	args := mock.Called(ctx)
	err := args.Error(1)

	if err != nil {
		return []dto.OrderResponse{}, err
	}

	return args.Get(0).([]dto.OrderResponse), nil
}

func (mock *MockGetOrdersWaitingPaymentUseCase) Execute(ctx context.Context) ([]dto.OrderResponse, error) {
	args := mock.Called(ctx)
	err := args.Error(1)

	if err != nil {
		return []dto.OrderResponse{}, err
	}

	return args.Get(0).([]dto.OrderResponse), nil
}

func (m *MockCreateOrderUseCase) Execute(
	ctx context.Context,
	order dto.Order,
	date int64,
	wg *sync.WaitGroup,
	_ chan bool) (dto.OrderResponse, error) {
	args := m.Called(ctx, order, date, wg, mock.Anything)
	err := args.Error(1)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return args.Get(0).(dto.OrderResponse), nil
}

func (m *MockCreateProductUseCase) Execute(ctx context.Context, product dto.ProductForm) (uint, error) {
	args := m.Called(ctx, product)
	err := args.Error(1)

	if err != nil {
		return uint(0), err
	}

	return args.Get(0).(uint), nil
}

func (m *MockGetProductsByCategoryUseCase) Execute(ctx context.Context, category string) ([]dto.ProductResponse, error) {
	args := m.Called(ctx, category)
	err := args.Error(1)

	if err != nil {
		return []dto.ProductResponse{}, err
	}

	return args.Get(0).([]dto.ProductResponse), nil
}

func (m *MockGetProductsByIDUseCase) Execute(ctx context.Context, id uint) (dto.ProductResponse, error) {
	args := m.Called(ctx, id)
	err := args.Error(1)

	if err != nil {
		return dto.ProductResponse{}, err
	}

	return args.Get(0).(dto.ProductResponse), nil
}

func (mock *MockCreateOrderUseCase) GenerateTicket(ctx context.Context, date int64) int {
	args := mock.Called(ctx, date)
	err := args.Error(1)

	if err != nil {
		return 0
	}

	return args.Get(0).(int)
}

func (mock *MockUpdateToPreparingUseCase) Execute(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(1)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockDeleteProductUseCase) Execute(ctx context.Context, productId uint) error {
	args := mock.Called(ctx, productId)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockUpdateProductUseCase) Execute(ctx context.Context, product dto.ProductForm) error {
	args := mock.Called(ctx, product)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockUpdateToDoneUseCase) Execute(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(1)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockUpdateToDeliveredUseCase) Execute(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(1)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockUpdateToNotDeliveredUseCase) Execute(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(1)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockGetOrderByIdUseCase) Execute(ctx context.Context, orderId uint) (dto.OrderResponse, error) {
	args := mock.Called(ctx, orderId)
	err := args.Error(1)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return args.Get(0).(dto.OrderResponse), nil
}

func (mock *MockGetPaymentTypesUseCase) Execute() []string {
	args := mock.Called()
	return args.Get(0).([]string)
}

func (mock *MockGetCategoryUseCase) Execute() []string {
	args := mock.Called()
	return args.Get(0).([]string)
}
