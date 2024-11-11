package usecases

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

var (
	orderCreation = dto.Order{
		TotalPrice:   12345,
		PaymentID:    uint(1),
		TicketNumber: 1,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: 1,
			},
			{
				ProductID: 2,
			},
		},
	}
	customerId = uint(1)

	orderCreationWithCustomer = dto.Order{
		TotalPrice:   12345,
		PaymentID:    uint(1),
		TicketNumber: 1,
		CustomerID:   &customerId,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: 1,
			},
			{
				ProductID: 2,
			},
		},
	}

	orderCreationResponse = dto.OrderResponse{
		OrderId:      1,
		OrderDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		TicketNumber: 1,
		OrderStatus:  "Criado",
		OrderProduct: []dto.OrderProductResponse{
			{
				ProductID:   1,
				ProductName: "ProductName 1",
			},
			{
				ProductID:   2,
				ProductName: "ProductName 2",
			},
		},
	}

	customerName                      = "Customer Name"
	orderWithCustomerCreationResponse = dto.OrderResponse{
		OrderId:      1,
		OrderDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
		TicketNumber: 1,
		OrderStatus:  "Criado",
		CustomerName: &customerName,
		OrderProduct: []dto.OrderProductResponse{
			{
				ProductID:   1,
				ProductName: "ProductName 1",
			},
			{
				ProductID:   2,
				ProductName: "ProductName 2",
			},
		},
	}

	ordersList = []dto.OrderResponse{
		{
			OrderId:      1,
			OrderDate:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			TicketNumber: 1,
			OrderStatus:  "Criado",
			OrderProduct: []dto.OrderProductResponse{
				{
					ProductID:   1,
					ProductName: "ProductName 1",
				},
				{
					ProductID:   2,
					ProductName: "ProductName 2",
				},
			},
		},
	}

	paymentCreation = dto.Payment{
		TotalPrice:  1234,
		PaymentType: "Crédito",
	}

	paymentResponse = dto.PaymentResponse{
		PaymentId:        1,
		PaymentGatewayId: "123",
		PaymentDate:      time.Date(2024, 10, 10, 0, 0, 0, 0, time.Local),
	}

	productCreation = dto.ProductForm{
		Name:        "Name",
		Description: "Description",
		Category:    "Category",
		Price:       23456,
		Images: []dto.ProducImage{
			{
				ImageUrl: "imageUrl",
			},
		},
	}

	productUpdate = dto.ProductForm{
		Id:          uint(12),
		Name:        "Name",
		Description: "Description",
		Category:    "Category",
		Price:       23456,
		Images: []dto.ProducImage{
			{
				ImageUrl: "imageUrl",
			},
		},
	}

	productsByCategory = []dto.ProductResponse{
		{
			Id:          uint(12),
			Name:        "Name",
			Description: "Description",
			Category:    "Category",
			Price:       23456,
			Images: []dto.ProducImage{
				{
					ImageUrl: "imageUrl",
				},
			},
		},
		{
			Id:          uint(23),
			Name:        "Name 2",
			Description: "Description 2",
			Category:    "Category 2",
			Price:       23456,
			Images: []dto.ProducImage{
				{
					ImageUrl: "imageUrl",
				},
			},
		},
		{
			Id:          uint(34),
			Name:        "Name 3",
			Description: "Description 3",
			Category:    "Category 3",
			Price:       34567,
			Images: []dto.ProducImage{
				{
					ImageUrl: "imageUrl",
				},
			},
		},
	}

	productById = dto.ProductResponse{
		Id:          uint(12),
		Name:        "Name",
		Description: "Description",
		Category:    "Category",
		Price:       23456,
		Images: []dto.ProducImage{
			{
				ImageUrl: "imageUrl",
			},
		},
	}
)

type MockOrderRepository struct {
	mock.Mock
}

type MockProductRepository struct {
	mock.Mock
}

type MockUserAdminRepository struct {
	mock.Mock
}

type MockQRCodePaymentRepository struct {
	mock.Mock
}

func (mock *MockOrderRepository) CreateOrder(ctx context.Context, order dto.Order) (dto.OrderResponse, error) {
	args := mock.Called(ctx, order)
	err := args.Error(1)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return args.Get(0).(dto.OrderResponse), nil
}

func (mock *MockOrderRepository) CreatePayingOrder(ctx context.Context, order dto.Order) (dto.OrderResponse, error) {
	args := mock.Called(ctx, order)
	err := args.Error(1)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return args.Get(0).(dto.OrderResponse), nil
}

func (mock *MockOrderRepository) DeleteOrder(ctx context.Context, orderID uint) error {
	args := mock.Called(ctx, orderID)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockOrderRepository) FinishOrderWithPayment(ctx context.Context, orderID uint, paymentID uint) error {
	args := mock.Called(ctx, orderID, paymentID)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockOrderRepository) GetOrderById(ctx context.Context, orderId uint) (dto.OrderResponse, error) {
	args := mock.Called(ctx, orderId)
	err := args.Error(1)

	if err != nil {
		return dto.OrderResponse{}, err
	}

	return args.Get(0).(dto.OrderResponse), nil
}

func (mock *MockOrderRepository) GetOrdersToPrepare(ctx context.Context) ([]dto.OrderResponse, error) {
	args := mock.Called(ctx)
	err := args.Error(1)

	if err != nil {
		return []dto.OrderResponse{}, err
	}

	return args.Get(0).([]dto.OrderResponse), nil
}

func (mock *MockOrderRepository) GetOrdersToFollow(ctx context.Context) ([]dto.OrderResponse, error) {
	args := mock.Called(ctx)
	err := args.Error(1)

	if err != nil {
		return []dto.OrderResponse{}, err
	}

	return args.Get(0).([]dto.OrderResponse), nil
}

func (mock *MockOrderRepository) GetOrdersWaitingPayment(ctx context.Context) ([]dto.OrderResponse, error) {
	args := mock.Called(ctx)
	err := args.Error(1)

	if err != nil {
		return []dto.OrderResponse{}, err
	}

	return args.Get(0).([]dto.OrderResponse), nil
}

func (mock *MockOrderRepository) UpdateToPreparing(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockOrderRepository) UpdateToDone(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockOrderRepository) UpdateToDelivered(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockOrderRepository) UpdateToNotDelivered(ctx context.Context, orderId uint) error {
	args := mock.Called(ctx, orderId)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockOrderRepository) GetNextTicketNumber(ctx context.Context, date int64) int {
	args := mock.Called(ctx, date)
	return args.Get(0).(int)
}

func (mock *MockProductRepository) CreateProduct(ctx context.Context, product dto.ProductForm) (uint, error) {
	args := mock.Called(ctx, product)
	err := args.Error(1)

	if err != nil {
		return 0, err
	}

	return args.Get(0).(uint), nil
}

func (mock *MockProductRepository) GetProductsByCategory(ctx context.Context, category string) ([]dto.ProductResponse, error) {
	args := mock.Called(ctx, category)
	err := args.Error(1)

	if err != nil {
		return []dto.ProductResponse{}, err
	}

	return args.Get(0).([]dto.ProductResponse), nil
}

func (mock *MockProductRepository) GetProductById(ctx context.Context, id uint) (dto.ProductResponse, error) {
	args := mock.Called(ctx, id)
	err := args.Error(1)

	if err != nil {
		return dto.ProductResponse{}, err
	}

	return args.Get(0).(dto.ProductResponse), nil
}

func (mock *MockProductRepository) DeleteProduct(ctx context.Context, productId uint) error {
	args := mock.Called(ctx, productId)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockProductRepository) UpdateProduct(ctx context.Context, product dto.ProductForm) error {
	args := mock.Called(ctx, product)
	err := args.Error(0)

	if err != nil {
		return err
	}

	return nil
}

func (mock *MockProductRepository) GetCategories() []string {
	args := mock.Called()
	return args.Get(0).([]string)
}
