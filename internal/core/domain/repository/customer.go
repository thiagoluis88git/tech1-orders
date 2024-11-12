package repository

import (
	"context"

	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

type CustomerRepository interface {
	GetCustomerByCPF(ctx context.Context, cpf string) (dto.Customer, error)
}
