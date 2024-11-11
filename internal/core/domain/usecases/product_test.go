package usecases

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

var (
	uc = NewValidateProductCategoryUseCase()
)

func TestProductsUseCase(t *testing.T) {
	t.Parallel()

	t.Run("got success when getting product categories in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewGetCategoriesUseCase(mockRepo)

		mockRepo.On("GetCategories").Return([]string{"Combo", "Bebidas", "Lanches"})

		response := sut.Execute()

		assert.NotEmpty(t, response)

		assert.Equal(t, 3, len(response))
	})

	t.Run("got success when getting products by category in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewGetProductsByCategoryUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetProductsByCategory", ctx, "category").Return(productsByCategory, nil)

		response, err := sut.Execute(ctx, "category")

		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)

		assert.Equal(t, 3, len(response))
	})

	t.Run("got error when getting products by category in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewGetProductsByCategoryUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetProductsByCategory", ctx, "category").Return(uint(0), &responses.LocalError{
			Code:    3,
			Message: "DATABASE_CONFLICT_ERROR",
		})

		response, err := sut.Execute(ctx, "category")

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusConflict, businessError.StatusCode)
	})

	t.Run("got success when getting product by id in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewGetProductByIdUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetProductById", ctx, uint(1)).Return(productById, nil)

		response, err := sut.Execute(ctx, uint(1))

		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)

		assert.Equal(t, uint(12), response.Id)
	})

	t.Run("got error when getting product by id in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewGetProductByIdUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("GetProductById", ctx, uint(1)).Return(dto.ProductResponse{}, &responses.LocalError{
			Code:    3,
			Message: "DATABASE_CONFLICT_ERROR",
		})

		response, err := sut.Execute(ctx, uint(1))

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusConflict, businessError.StatusCode)
	})

	t.Run("got success when creating product in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewCreateProductUseCase(uc, mockRepo)

		ctx := context.TODO()

		mockRepo.On("CreateProduct", ctx, productCreation).Return(uint(1), nil)

		response, err := sut.Execute(ctx, productCreation)

		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)

		assert.Equal(t, uint(1), response)
	})

	t.Run("got error when creating product in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewCreateProductUseCase(uc, mockRepo)

		ctx := context.TODO()

		mockRepo.On("CreateProduct", ctx, productCreation).Return(uint(0), &responses.LocalError{
			Code:    3,
			Message: "DATABASE_CONFLICT_ERROR",
		})

		response, err := sut.Execute(ctx, productCreation)

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Empty(t, response)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusConflict, businessError.StatusCode)
	})

	t.Run("got success when deleting product in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewDeleteProductUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("DeleteProduct", ctx, uint(12)).Return(nil)

		err := sut.Execute(ctx, uint(12))

		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
	})

	t.Run("got error when deleting product in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewDeleteProductUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("DeleteProduct", ctx, uint(12)).Return(&responses.LocalError{
			Code:    3,
			Message: "DATABASE_CONFLICT_ERROR",
		})

		err := sut.Execute(ctx, uint(12))

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusConflict, businessError.StatusCode)
	})

	t.Run("got success when updating product in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewUpdateProductUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("UpdateProduct", ctx, productUpdate).Return(nil)

		err := sut.Execute(ctx, productUpdate)

		mockRepo.AssertExpectations(t)

		assert.NoError(t, err)
	})

	t.Run("got error when updating product in services", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockProductRepository)
		sut := NewUpdateProductUseCase(mockRepo)

		ctx := context.TODO()

		mockRepo.On("UpdateProduct", ctx, productUpdate).Return(&responses.LocalError{
			Code:    3,
			Message: "DATABASE_CONFLICT_ERROR",
		})

		err := sut.Execute(ctx, productUpdate)

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)

		var businessError *responses.BusinessResponse
		assert.Equal(t, true, errors.As(err, &businessError))
		assert.Equal(t, http.StatusConflict, businessError.StatusCode)
	})
}
