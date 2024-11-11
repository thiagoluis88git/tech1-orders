package usecases

import (
	"context"
	"slices"

	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/repository"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

type ValidateOrderToDeliveredOrNotUseCase struct {
	repository repository.OrderRepository
}

type SortOrdersUseCase struct{}

type ValidateOrderToDoneUseCase struct {
	repository repository.OrderRepository
}

type ValidateOrderToPrepareUseCase struct {
	repository repository.OrderRepository
}

func NewValidateOrderToPrepareUseCase(repository repository.OrderRepository) *ValidateOrderToPrepareUseCase {
	return &ValidateOrderToPrepareUseCase{
		repository: repository,
	}
}

func NewValidateOrderToDoneUseCase(repository repository.OrderRepository) *ValidateOrderToDoneUseCase {
	return &ValidateOrderToDoneUseCase{
		repository: repository,
	}
}

func NewSortOrdersUseCase() *SortOrdersUseCase {
	return &SortOrdersUseCase{}
}

func NewValidateOrderToDeliveredOrNotUseCase(repository repository.OrderRepository) *ValidateOrderToDeliveredOrNotUseCase {
	return &ValidateOrderToDeliveredOrNotUseCase{
		repository: repository,
	}
}

func (usecase *ValidateOrderToDeliveredOrNotUseCase) Execute(ctx context.Context, orderId uint) error {
	response, err := usecase.repository.GetOrderById(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "ValidateOrderToDoneUseCase -> GetOrderById")
	}

	if response.OrderStatus != "Finalizado" {
		return &responses.BusinessResponse{
			StatusCode: 428,
			Message:    "The order must be in Finalizado status",
		}
	}

	return nil
}

func (usecase *SortOrdersUseCase) Execute(orders []dto.OrderResponse) {
	slices.SortFunc(orders, func(previous, next dto.OrderResponse) int {
		if next.OrderStatus == "Finalizado" && (previous.OrderStatus == "Preparando" || previous.OrderStatus == "Criado") {
			return 1
		}

		if next.OrderStatus == "Preparando" && previous.OrderStatus == "Criado" {
			return 0
		}

		return -1
	})
}

func (usecase *ValidateOrderToDoneUseCase) Execute(ctx context.Context, orderId uint) error {
	response, err := usecase.repository.GetOrderById(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "ValidateOrderToDoneUseCase -> GetOrderById")
	}

	if response.OrderStatus != "Preparando" {
		return &responses.BusinessResponse{
			StatusCode: 428,
			Message:    "The order must be in Preparando status",
		}
	}

	return nil
}

func (usecase *ValidateOrderToPrepareUseCase) Execute(ctx context.Context, orderId uint) error {
	response, err := usecase.repository.GetOrderById(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "ValidateOrderToPrepareUseCase -> GetOrderById")
	}

	if response.OrderStatus != "Criado" {
		return &responses.BusinessResponse{
			StatusCode: 428,
			Message:    "The order must be in Criado status",
		}
	}

	return nil
}
