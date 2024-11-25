package bdd_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

type MockCreateProductUseCase struct {
	mock.Mock
}

func (m *MockCreateProductUseCase) Execute(ctx context.Context, product dto.ProductForm) (uint, error) {
	args := m.Called(ctx, product)
	err := args.Error(1)

	if err != nil {
		return uint(0), err
	}

	return args.Get(0).(uint), nil
}
