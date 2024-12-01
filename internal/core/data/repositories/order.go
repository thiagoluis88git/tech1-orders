package repositories

import (
	"context"
	"time"

	"github.com/thiagoluis88git/tech1-orders/internal/core/data/model"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/remote"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/repository"
	"github.com/thiagoluis88git/tech1-orders/pkg/database"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

type OrderRespository struct {
	db         *database.Database
	customerDS remote.CustomerRemoteDataSource
}

func NewOrderRespository(db *database.Database, customerDS remote.CustomerRemoteDataSource) repository.OrderRepository {
	return &OrderRespository{
		db:         db,
		customerDS: customerDS,
	}
}

func (repository *OrderRespository) CreateOrder(ctx context.Context, order dto.Order) (dto.OrderResponse, error) {
	return repository.createOrder(ctx, order, model.OrderStatusCreated)
}

func (repository *OrderRespository) CreatePayingOrder(ctx context.Context, order dto.Order) (dto.OrderResponse, error) {
	return repository.createOrder(ctx, order, model.OrderStatusPaying)
}

func (repository *OrderRespository) createOrder(ctx context.Context, order dto.Order, status string) (dto.OrderResponse, error) {
	tx := repository.db.Connection.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	orderEntity := &model.Order{
		OrderStatus:  status,
		TotalPrice:   order.TotalPrice,
		CPF:          order.CPF,
		PaymentID:    order.PaymentID,
		TicketNumber: order.TicketNumber,
	}

	err := tx.Create(orderEntity).Error

	if err != nil {
		tx.Rollback()
		return dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	orderProductsEntity := []*model.OrderProduct{}

	for _, value := range order.OrderProduct {
		orderProductsEntity = append(orderProductsEntity, &model.OrderProduct{
			ProductID: value.ProductID,
			OrderID:   orderEntity.ID,
		})
	}

	err = tx.Create(orderProductsEntity).Error

	if err != nil {
		tx.Rollback()
		return dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	err = tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	return dto.OrderResponse{
		OrderId:      orderEntity.ID,
		OrderDate:    orderEntity.CreatedAt,
		TicketNumber: orderEntity.TicketNumber,
	}, nil
}

func (repository *OrderRespository) DeleteOrder(ctx context.Context, orderID uint) error {
	tx := repository.db.Connection.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return responses.GetDatabaseError(err)
	}

	err := tx.Where("order_id = ?", orderID).Delete(&model.OrderProduct{}).Error

	if err != nil {
		tx.Rollback()
		return responses.GetDatabaseError(err)
	}

	err = tx.Delete(&model.Order{}, orderID).Error

	if err != nil {
		tx.Rollback()
		return responses.GetDatabaseError(err)
	}

	err = tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *OrderRespository) FinishOrderWithPayment(ctx context.Context, orderID uint, paymentID string) error {
	err := repository.db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("payment_id", paymentID).
		Update("order_status", model.OrderStatusCreated).
		Error

	if err != nil {
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *OrderRespository) GetOrderById(ctx context.Context, orderId uint) (dto.OrderResponse, error) {
	var orderEntity model.Order
	err := repository.
		db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Preload("OrderProduct.Product").
		Where("id = ?", orderId).
		Find(&orderEntity).
		Limit(1).
		Error

	if err != nil {
		return dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	if orderEntity.ID == uint(0) {
		return dto.OrderResponse{}, &responses.LocalError{
			Message: "Order not found",
			Code:    responses.NOT_FOUND_ERROR,
		}
	}

	orderProduct := []dto.OrderProductResponse{}

	for _, value := range orderEntity.OrderProduct {
		orderProduct = append(orderProduct, dto.OrderProductResponse{
			ProductID:   value.ProductID,
			ProductName: value.Product.Name,
			Description: value.Product.Description,
		})
	}

	var customerName *string

	if orderEntity.CPF != nil {
		customerResponse, err := repository.customerDS.GetCustomerByCPF(ctx, *orderEntity.CPF)

		if err != nil {
			return dto.OrderResponse{}, err
		}

		customerName = &customerResponse.Name
	}

	return dto.OrderResponse{
		OrderId:        orderEntity.ID,
		OrderDate:      orderEntity.CreatedAt,
		PreparingAt:    orderEntity.PreparingAt,
		DoneAt:         orderEntity.DoneAt,
		DeliveredAt:    orderEntity.DeliveredAt,
		NotDeliveredAt: orderEntity.NotDeliveredAt,
		TicketNumber:   orderEntity.TicketNumber,
		OrderStatus:    orderEntity.OrderStatus,
		OrderProduct:   orderProduct,
		CustomerName:   customerName,
	}, nil
}

func (repository *OrderRespository) GetOrdersToPrepare(ctx context.Context) ([]dto.OrderResponse, error) {
	var orderEntity []model.Order
	err := repository.
		db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Preload("OrderProduct.Product").
		Where("order_status = ?", model.OrderStatusCreated).
		Order("created_at").
		Find(&orderEntity).
		Error

	if err != nil {
		return []dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	return repository.buildOrdersList(orderEntity), nil
}

func (repository *OrderRespository) GetOrdersToFollow(ctx context.Context) ([]dto.OrderResponse, error) {
	var orderEntity []model.Order
	err := repository.
		db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Preload("OrderProduct.Product").
		Where("order_status in (?, ?,?)",
			model.OrderStatusCreated,
			model.OrderStatusPreparing,
			model.OrderStatusDone,
		).
		Order("created_at").
		Find(&orderEntity).
		Error

	if err != nil {
		return []dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	return repository.buildOrdersList(orderEntity), nil
}

func (repository *OrderRespository) GetOrdersWaitingPayment(ctx context.Context) ([]dto.OrderResponse, error) {
	var orderEntity []model.Order
	err := repository.
		db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Preload("OrderProduct.Product").
		Where("order_status = ?", model.OrderStatusPaying).
		Order("created_at").
		Find(&orderEntity).
		Error

	if err != nil {
		return []dto.OrderResponse{}, responses.GetDatabaseError(err)
	}

	return repository.buildOrdersList(orderEntity), nil
}

func (repository *OrderRespository) buildOrdersList(orderEntity []model.Order) []dto.OrderResponse {
	orders := []dto.OrderResponse{}

	for _, value := range orderEntity {
		orderProduct := []dto.OrderProductResponse{}

		for _, value := range value.OrderProduct {
			orderProduct = append(orderProduct, dto.OrderProductResponse{
				ProductID:   value.ProductID,
				ProductName: value.Product.Name,
				Description: value.Product.Description,
			})
		}

		var customerName *string

		if value.CPF != nil {
			customerResponse, err := repository.customerDS.GetCustomerByCPF(context.Background(), *value.CPF)

			if err == nil {
				customerName = &customerResponse.Name
			}
		}

		orders = append(orders, dto.OrderResponse{
			OrderId:        value.ID,
			OrderDate:      value.CreatedAt,
			PreparingAt:    value.PreparingAt,
			DoneAt:         value.DoneAt,
			DeliveredAt:    value.DeliveredAt,
			NotDeliveredAt: value.NotDeliveredAt,
			TicketNumber:   value.TicketNumber,
			OrderStatus:    value.OrderStatus,
			OrderProduct:   orderProduct,
			CustomerName:   customerName,
		})
	}

	return orders
}

func (repository *OrderRespository) UpdateToPreparing(ctx context.Context, orderId uint) error {
	err := repository.db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderId).
		Update("order_status", model.OrderStatusPreparing).
		Update("preparing_at", time.Now()).
		Error

	if err != nil {
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *OrderRespository) UpdateToDone(ctx context.Context, orderId uint) error {
	err := repository.db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderId).
		Update("order_status", model.OrderStatusDone).
		Update("done_at", time.Now()).
		Error

	if err != nil {
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *OrderRespository) UpdateToDelivered(ctx context.Context, orderId uint) error {
	err := repository.db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderId).
		Update("order_status", model.OrderStatusDelivered).
		Update("delivered_at", time.Now()).
		Error

	if err != nil {
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *OrderRespository) UpdateToNotDelivered(ctx context.Context, orderId uint) error {
	err := repository.db.Connection.WithContext(ctx).
		Model(&model.Order{}).
		Where("id = ?", orderId).
		Update("order_status", model.OrderStatusNotDelivered).
		Update("not_delivered_at", time.Now()).
		Error

	if err != nil {
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *OrderRespository) GetNextTicketNumber(ctx context.Context, date int64) int {
	var orderTicketNumber model.OrderTicketNumber
	err := repository.db.Connection.WithContext(ctx).
		Model(&model.OrderTicketNumber{}).
		Where("date = ?", date).
		Find(&orderTicketNumber).
		Limit(1).
		Error

	if err != nil || orderTicketNumber.Date == 0 {
		return repository.createNewTicketForDate(ctx, date)
	}

	newTicketNumber := orderTicketNumber.TicketNumber + 1

	return repository.updateTicketForDate(ctx, date, newTicketNumber)
}

func (repository *OrderRespository) createNewTicketForDate(ctx context.Context, date int64) int {
	orderTicketNumber := model.OrderTicketNumber{
		Date:         date,
		TicketNumber: 1,
	}

	errCreate := repository.db.Connection.WithContext(ctx).Create(&orderTicketNumber).Error

	if errCreate != nil {
		return 999
	}

	return 1
}

func (repository *OrderRespository) updateTicketForDate(ctx context.Context, date int64, newTicketNumber int) int {
	err := repository.db.Connection.WithContext(ctx).
		Model(&model.OrderTicketNumber{}).
		Where("date = ?", date).
		Update("ticket_number", newTicketNumber).
		Error

	if err != nil {
		return 999
	}

	return newTicketNumber
}
