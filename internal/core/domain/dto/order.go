package dto

import "time"

type Order struct {
	OrderStatus  string
	TotalPrice   float64        `json:"totalPrice" validate:"required"`
	CustomerID   *uint          `json:"customerId"`
	PaymentID    uint           `json:"paymentId" validate:"required"`
	OrderProduct []OrderProduct `json:"orderProducts" validate:"required"`
	TicketNumber int
}

type QRCodeOrder struct {
	OrderStatus  string
	TotalPrice   float64        `json:"totalPrice" validate:"required"`
	CustomerID   *uint          `json:"customerId"`
	OrderProduct []OrderProduct `json:"orderProducts" validate:"required"`
	TicketNumber int
	PaymentID    uint
}

type OrderProduct struct {
	ProductID    uint    `json:"productId" validate:"required"`
	ProductPrice float64 `json:"productPrice" validate:"required"`
}

type OrderResponse struct {
	OrderId        uint                   `json:"orderId"`
	OrderDate      time.Time              `json:"orderDate"`
	PreparingAt    *time.Time             `json:"preparingAt"`
	DoneAt         *time.Time             `json:"doneAt"`
	DeliveredAt    *time.Time             `json:"deliveredAt"`
	NotDeliveredAt *time.Time             `json:"notDeliveredAt"`
	TicketNumber   int                    `json:"ticketNumber"`
	CustomerName   *string                `json:"customerName"`
	OrderStatus    string                 `json:"orderStatus"`
	OrderProduct   []OrderProductResponse `json:"orderProducts"`
}

type OrderProductResponse struct {
	ProductID   uint   `json:"id"`
	ProductName string `json:"name"`
	Description string `json:"description"`
}
