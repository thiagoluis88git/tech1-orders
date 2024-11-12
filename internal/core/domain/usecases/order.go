package usecases

import (
	"context"
	"sync"

	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/repository"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

type CreateOrderUseCase interface {
	Execute(ctx context.Context, order dto.Order, date int64, wg *sync.WaitGroup, ch chan bool) (dto.OrderResponse, error)
	GenerateTicket(ctx context.Context, date int64) int
}

type CreateOrderUseCaseImpl struct {
	orderRepo                repository.OrderRepository
	customerRepo             repository.CustomerRepository
	validateToPrepare        *ValidateOrderToPrepareUseCase
	validateToDone           *ValidateOrderToDoneUseCase
	validateToDeliveredOrNot *ValidateOrderToDeliveredOrNotUseCase
	sortOrderUseCase         *SortOrdersUseCase
}

type UpdateToPreparingUseCase interface {
	Execute(ctx context.Context, orderId uint) error
}

type UpdateToPreparingUseCaseImpl struct {
	orderRepo         repository.OrderRepository
	validateToPrepare *ValidateOrderToPrepareUseCase
}

type UpdateToDoneUseCase interface {
	Execute(ctx context.Context, orderId uint) error
}

type UpdateToDoneUseCaseImpl struct {
	orderRepo      repository.OrderRepository
	validateToDone *ValidateOrderToDoneUseCase
}

type UpdateToDeliveredUseCase interface {
	Execute(ctx context.Context, orderId uint) error
}

type UpdateToDeliveredUseCaseImpl struct {
	orderRepo                repository.OrderRepository
	validateToDeliveredOrNot *ValidateOrderToDeliveredOrNotUseCase
}

type UpdateToNotDeliveredUseCase interface {
	Execute(ctx context.Context, orderId uint) error
}

type UpdateToNotDeliveredUseCaseImpl struct {
	orderRepo                repository.OrderRepository
	validateToDeliveredOrNot *ValidateOrderToDeliveredOrNotUseCase
}

type GetOrderByIdUseCase interface {
	Execute(ctx context.Context, orderId uint) (dto.OrderResponse, error)
}

type GetOrderByIdUseCaseImpl struct {
	orderRepo repository.OrderRepository
}

type GetOrdersToPrepareUseCase interface {
	Execute(ctx context.Context) ([]dto.OrderResponse, error)
}

type GetOrdersToPrepareUseCaseImpl struct {
	orderRepo        repository.OrderRepository
	sortOrderUseCase *SortOrdersUseCase
}

type GetOrdersToFollowUseCase interface {
	Execute(ctx context.Context) ([]dto.OrderResponse, error)
}

type GetOrdersToFollowUseCaseImpl struct {
	orderRepo        repository.OrderRepository
	sortOrderUseCase *SortOrdersUseCase
}

type GetOrdersWaitingPaymentUseCase interface {
	Execute(ctx context.Context) ([]dto.OrderResponse, error)
}

type GetOrdersWaitingPaymentUseCaseImpl struct {
	orderRepo        repository.OrderRepository
	sortOrderUseCase *SortOrdersUseCase
}

func NewCreateOrderUseCase(
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	validateToPrepate *ValidateOrderToPrepareUseCase,
	validateToDone *ValidateOrderToDoneUseCase,
	validateToDeliveredOrNot *ValidateOrderToDeliveredOrNotUseCase,
	sortOrderUseCase *SortOrdersUseCase,
) CreateOrderUseCase {
	return &CreateOrderUseCaseImpl{
		orderRepo:                orderRepo,
		customerRepo:             customerRepo,
		validateToPrepare:        validateToPrepate,
		validateToDone:           validateToDone,
		validateToDeliveredOrNot: validateToDeliveredOrNot,
		sortOrderUseCase:         sortOrderUseCase,
	}
}

func NewGetOrderByIdUseCase(
	orderRepo repository.OrderRepository,
) GetOrderByIdUseCase {
	return &GetOrderByIdUseCaseImpl{
		orderRepo: orderRepo,
	}
}

func NewGetOrdersToPrepareUseCase(
	orderRepo repository.OrderRepository,
	sortOrderUseCase *SortOrdersUseCase,
) GetOrdersToPrepareUseCase {
	return &GetOrdersToPrepareUseCaseImpl{
		orderRepo:        orderRepo,
		sortOrderUseCase: sortOrderUseCase,
	}
}

func NewGetOrdersToFollowUseCase(
	orderRepo repository.OrderRepository,
	sortOrderUseCase *SortOrdersUseCase,
) GetOrdersToFollowUseCase {
	return &GetOrdersToFollowUseCaseImpl{
		orderRepo:        orderRepo,
		sortOrderUseCase: sortOrderUseCase,
	}
}

func NewGetOrdersWaitingPaymentUseCase(
	orderRepo repository.OrderRepository,
	sortOrderUseCase *SortOrdersUseCase,
) GetOrdersWaitingPaymentUseCase {
	return &GetOrdersWaitingPaymentUseCaseImpl{
		orderRepo:        orderRepo,
		sortOrderUseCase: sortOrderUseCase,
	}
}

func NewUpdateToPreparingUseCase(
	orderRepo repository.OrderRepository,
	validateToPrepare *ValidateOrderToPrepareUseCase,
) UpdateToPreparingUseCase {
	return &UpdateToPreparingUseCaseImpl{
		orderRepo:         orderRepo,
		validateToPrepare: validateToPrepare,
	}
}

func NewUpdateToDoneUseCase(
	orderRepo repository.OrderRepository,
	validateToDone *ValidateOrderToDoneUseCase,
) UpdateToDoneUseCase {
	return &UpdateToDoneUseCaseImpl{
		orderRepo:      orderRepo,
		validateToDone: validateToDone,
	}
}

func NewUpdateToDeliveredUseCase(
	orderRepo repository.OrderRepository,
	validateToDeliveredOrNot *ValidateOrderToDeliveredOrNotUseCase,
) UpdateToDeliveredUseCase {
	return &UpdateToDeliveredUseCaseImpl{
		orderRepo:                orderRepo,
		validateToDeliveredOrNot: validateToDeliveredOrNot,
	}
}

func NewUpdateToNotDeliveredUseCase(
	orderRepo repository.OrderRepository,
	validateToDeliveredOrNot *ValidateOrderToDeliveredOrNotUseCase,
) UpdateToNotDeliveredUseCase {
	return &UpdateToNotDeliveredUseCaseImpl{
		orderRepo:                orderRepo,
		validateToDeliveredOrNot: validateToDeliveredOrNot,
	}
}

func (usecase *CreateOrderUseCaseImpl) Execute(
	ctx context.Context,
	order dto.Order,
	date int64,
	wg *sync.WaitGroup,
	ch chan bool,
) (dto.OrderResponse, error) {
	//Block this code below until this Channel be empty (by reading with <-ch)
	ch <- true

	order.TicketNumber = usecase.GenerateTicket(ctx, date)

	response, err := usecase.orderRepo.CreateOrder(ctx, order)

	if err != nil {
		return dto.OrderResponse{}, responses.GetResponseError(err, "OrderService -> CreateOrder")
	}

	if order.CPF != nil {
		customer, err := usecase.customerRepo.GetCustomerByCPF(ctx, *order.CPF)
		if err == nil {
			response.CustomerName = &customer.Name
		}
	}

	// Release the channel to others process be able to start a new order creation
	<-ch
	wg.Done()

	return response, nil
}

func (usecase *CreateOrderUseCaseImpl) GenerateTicket(ctx context.Context, date int64) int {
	return usecase.orderRepo.GetNextTicketNumber(ctx, date)
}

func (usecase *GetOrderByIdUseCaseImpl) Execute(ctx context.Context, orderId uint) (dto.OrderResponse, error) {
	response, err := usecase.orderRepo.GetOrderById(ctx, orderId)

	if err != nil {
		return dto.OrderResponse{}, responses.GetResponseError(err, "OrderService -> GetOrderById")
	}

	return response, nil
}

func (usecase *GetOrdersToPrepareUseCaseImpl) Execute(ctx context.Context) ([]dto.OrderResponse, error) {
	response, err := usecase.orderRepo.GetOrdersToPrepare(ctx)

	if err != nil {
		return []dto.OrderResponse{}, responses.GetResponseError(err, "OrderService -> GetOrdersToPrepare")
	}

	usecase.sortOrderUseCase.Execute(response)

	return response, nil
}

func (usecase *GetOrdersToFollowUseCaseImpl) Execute(ctx context.Context) ([]dto.OrderResponse, error) {
	response, err := usecase.orderRepo.GetOrdersToFollow(ctx)

	if err != nil {
		return []dto.OrderResponse{}, responses.GetResponseError(err, "OrderService -> GetOrdersToFollow")
	}

	usecase.sortOrderUseCase.Execute(response)

	return response, nil
}

func (usecase *GetOrdersWaitingPaymentUseCaseImpl) Execute(ctx context.Context) ([]dto.OrderResponse, error) {
	response, err := usecase.orderRepo.GetOrdersWaitingPayment(ctx)

	if err != nil {
		return []dto.OrderResponse{}, responses.GetResponseError(err, "OrderService -> GetOrdersWaitingPayment")
	}

	usecase.sortOrderUseCase.Execute(response)

	return response, nil
}

func (usecase *UpdateToPreparingUseCaseImpl) Execute(ctx context.Context, orderId uint) error {
	err := usecase.validateToPrepare.Execute(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToPreparing")
	}

	err = usecase.orderRepo.UpdateToPreparing(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToPreparing")
	}

	return nil
}

func (usecase *UpdateToDoneUseCaseImpl) Execute(ctx context.Context, orderId uint) error {
	err := usecase.validateToDone.Execute(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToDone")
	}

	err = usecase.orderRepo.UpdateToDone(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToDone")
	}

	return nil
}

func (usecase *UpdateToDeliveredUseCaseImpl) Execute(ctx context.Context, orderId uint) error {
	err := usecase.validateToDeliveredOrNot.Execute(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToDelivered")
	}

	err = usecase.orderRepo.UpdateToDelivered(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToDelivered")
	}

	return nil
}

func (usecase *UpdateToNotDeliveredUseCaseImpl) Execute(ctx context.Context, orderId uint) error {
	err := usecase.validateToDeliveredOrNot.Execute(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToNotDelivered")
	}

	err = usecase.orderRepo.UpdateToNotDelivered(ctx, orderId)

	if err != nil {
		return responses.GetResponseError(err, "OrderService -> UpdateToNotDelivered")
	}

	return nil
}
