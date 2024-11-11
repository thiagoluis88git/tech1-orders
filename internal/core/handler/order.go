package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/usecases"
	"github.com/thiagoluis88git/tech1-orders/pkg/httpserver"
)

// @Summary Create new order
// @Description Create new order. To make an order the payment needs to be completed
// @Description A new Ticket will be generated by the Order Date starting from 1
// @Description In the next day the Ticket number will starts from 1 and so on
// @Tags Order
// @Accept json
// @Produce json
// @Param product body dto.Order true "order"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 "Order has required fields"
// @Router /api/orders [post]
func CreateOrderHandler(createOrder usecases.CreateOrderUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order dto.Order

		err := httpserver.DecodeJSONBody(w, r, &order)

		if err != nil {
			log.Print("decoding order body", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		now := time.Now()
		orderDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		// Create this to prevent 2 process/goroutines create order with the same TicketNumber
		var waitGroup sync.WaitGroup
		ch := make(chan bool, 1)
		waitGroup.Add(1)

		response, err := createOrder.Execute(r.Context(), order, orderDate.UnixMilli(), &waitGroup, ch)

		if err != nil {
			log.Print("create order", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, response)
	}
}

// @Summary Get order by Id
// @Description Get an order by Id
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "12"
// @Success 200 {object} dto.OrderResponse
// @Failure 400 "Order has required fields"
// @Router /api/orders/{id} [get]
func GetOrderByIdHandler(getOrderById usecases.GetOrderByIdUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderIdStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("get order by id path", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		orderId, err := strconv.Atoi(orderIdStr)

		if err != nil {
			log.Print("get order by id path", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		response, err := getOrderById.Execute(r.Context(), uint(orderId))

		if err != nil {
			log.Print("get order by id", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, response)
	}
}

// @Summary Get all orders to prepare
// @Description Get all orders already payed that needs to be prepared. This endpoint will be used by the kitchen
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} []dto.OrderResponse
// @Router /api/orders/to-prepare [get]
func GetOrdersToPrepareHandler(getOrdersToPrepare usecases.GetOrdersToPrepareUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := getOrdersToPrepare.Execute(r.Context())

		if err != nil {
			log.Print("get orders to prepare", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, response)
	}
}

// @Summary Get all orders status different to prepare
// @Description Get all orders status by the waiter and the customer. This endpoint will be used by the waiter and customer
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} []dto.OrderResponse
// @Router /api/orders/status [get]
func GetOrdersToFollowHandler(getOrdersToFollow usecases.GetOrdersToFollowUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := getOrdersToFollow.Execute(r.Context())

		if err != nil {
			log.Print("get orders status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, response)
	}
}

// @Summary Get all orders with waiting payment status
// @Description Get all orders with waiting payment by the owner.
// @Description This endpoint will be used by the owner to know it the Mercado Livre QR Code was paid
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} []dto.OrderResponse
// @Router /api/orders/waiting-payment [get]
func GetOrdersWaitingPaymentHandler(getOrdersWaitingPayment usecases.GetOrdersWaitingPaymentUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := getOrdersWaitingPayment.Execute(r.Context())

		if err != nil {
			log.Print("get orders status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, response)
	}
}

// @Summary Update an order to PREPARING
// @Description Update an order. This service wil be used by the kitchen to notify a customer that the order is being prepared
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "12"
// @Success 204
// @Failure 404 "Order not found"
// @Failure 428 "Precondition failed: Need to be with status Criado"
// @Router /api/orders/{id}/preparing [put]
func UpdateOrderPreparingHandler(updateToPreparing usecases.UpdateToPreparingUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		id, err := getOrderId(idStr)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		err = updateToPreparing.Execute(r.Context(), id)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}

// @Summary Update an order to DONE
// @Description Update an order. This service wil be used by the kitchen to notify a customer and the waiter that the order is done
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "12"
// @Success 204
// @Failure 404 "Order not found"
// @Failure 428 "Precondition failed: Need to be with status Preparando"
// @Router /api/orders/{id}/done [put]
func UpdateOrderDoneHandler(updateToDone usecases.UpdateToDoneUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		id, err := getOrderId(idStr)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		err = updateToDone.Execute(r.Context(), id)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}

// @Summary Update an order to DELIVERED
// @Description Update an order. This service wil be used by the waiter to close the order informing that user got its order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "12"
// @Success 204
// @Failure 404 "Order not found"
// @Failure 428 "Precondition failed: Need to be with status Finalizado"
// @Router /api/orders/{id}/delivered [put]
func UpdateOrderDeliveredHandler(updateToDelivered usecases.UpdateToDeliveredUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		id, err := getOrderId(idStr)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		err = updateToDelivered.Execute(r.Context(), id)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}

// @Summary Update an order to NOT_DELIVERED
// @Description Update an order. This service wil be used by the waiter to close the order informing that user didn't get the order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "12"
// @Success 204
// @Failure 404 "Order not found"
// @Failure 428 "Precondition failed: Need to be with status Finalizado"
// @Router /api/orders/{id}/not-delivered [put]
func UpdateOrderNotDeliveredandler(updateToNotDelivered usecases.UpdateToNotDeliveredUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		id, err := getOrderId(idStr)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		err = updateToNotDelivered.Execute(r.Context(), id)

		if err != nil {
			log.Print("update order status", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}

func getOrderId(orderdStr string) (uint, error) {
	orderId, err := strconv.Atoi(orderdStr)

	if err != nil {
		return 0, err
	}

	return uint(orderId), nil
}
