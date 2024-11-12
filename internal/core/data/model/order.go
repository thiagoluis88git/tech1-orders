package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	OrderStatusPaying       = "Em pagamento"
	OrderStatusCreated      = "Criado"
	OrderStatusPreparing    = "Preparando"
	OrderStatusDone         = "Finalizado"
	OrderStatusDelivered    = "Entregue"
	OrderStatusNotDelivered = "Não entregue"
)

type Order struct {
	gorm.Model
	OrderStatus    string
	TotalPrice     float64
	PaymentID      uint
	CPF            *string
	TicketNumber   int
	PreparingAt    *time.Time
	DoneAt         *time.Time
	DeliveredAt    *time.Time
	NotDeliveredAt *time.Time
	OrderProduct   []OrderProduct
}

type OrderProduct struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Product   Product
}

type OrderTicketNumber struct {
	Date         int64 `gorm:"index;unique"`
	TicketNumber int
}
