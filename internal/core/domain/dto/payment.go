package dto

import "time"

type Payment struct {
	CustomerID  *uint   `json:"customerId"`
	TotalPrice  float64 `json:"totalPrice" validate:"required"`
	PaymentType string  `json:"paymentType" validate:"required"`
}

type PaymentResponse struct {
	PaymentId        uint      `json:"paymentId"`
	PaymentGatewayId string    `json:"paymentGatewayId"`
	PaymentDate      time.Time `json:"paymentDate"`
}
