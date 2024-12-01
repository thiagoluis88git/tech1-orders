package repository

import (
	"context"

	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order dto.Order) (dto.OrderResponse, error)
	CreatePayingOrder(ctx context.Context, order dto.Order) (dto.OrderResponse, error)
	FinishOrderWithPayment(ctx context.Context, orderID uint, paymentID string) error
	DeleteOrder(ctx context.Context, orderID uint) error
	GetOrderById(ctx context.Context, orderID uint) (dto.OrderResponse, error)
	GetOrdersToPrepare(ctx context.Context) ([]dto.OrderResponse, error)
	GetOrdersToFollow(ctx context.Context) ([]dto.OrderResponse, error)
	GetOrdersWaitingPayment(ctx context.Context) ([]dto.OrderResponse, error)
	UpdateToPreparing(ctx context.Context, orderID uint) error
	UpdateToDone(ctx context.Context, orderID uint) error
	UpdateToDelivered(ctx context.Context, orderID uint) error
	UpdateToNotDelivered(ctx context.Context, orderID uint) error
	GetNextTicketNumber(ctx context.Context, date int64) int
}
