package repository

import (
	"context"

	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product dto.ProductForm) (uint, error)
	GetCategories() []string
	GetProductsByCategory(ctx context.Context, category string) ([]dto.ProductResponse, error)
	GetProductById(ctx context.Context, id uint) (dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, productId uint) error
	UpdateProduct(ctx context.Context, product dto.ProductForm) error
}
