package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

func TestSortOrdersUseCase(t *testing.T) {
	t.Parallel()

	t.Run("got success when sorting orders with many orders use case", func(t *testing.T) {
		t.Parallel()

		sut := NewSortOrdersUseCase()

		orders := []dto.OrderResponse{
			{
				OrderStatus: "Finalizado",
			},
			{
				OrderStatus: "Criado",
			},
			{
				OrderStatus: "Criado",
			},
			{
				OrderStatus: "Preparando",
			},
			{
				OrderStatus: "Preparando",
			},
			{
				OrderStatus: "Finalizado",
			},
			{
				OrderStatus: "Criado",
			},
			{
				OrderStatus: "Preparando",
			},
		}

		assert.Equal(t, "Finalizado", orders[0].OrderStatus)
		assert.Equal(t, "Criado", orders[1].OrderStatus)
		assert.Equal(t, "Criado", orders[2].OrderStatus)
		assert.Equal(t, "Preparando", orders[3].OrderStatus)
		assert.Equal(t, "Preparando", orders[4].OrderStatus)
		assert.Equal(t, "Finalizado", orders[5].OrderStatus)
		assert.Equal(t, "Criado", orders[6].OrderStatus)
		assert.Equal(t, "Preparando", orders[7].OrderStatus)

		sut.Execute(orders)

		assert.Equal(t, "Finalizado", orders[0].OrderStatus)
		assert.Equal(t, "Finalizado", orders[1].OrderStatus)
		assert.Equal(t, "Preparando", orders[2].OrderStatus)
		assert.Equal(t, "Preparando", orders[3].OrderStatus)
		assert.Equal(t, "Preparando", orders[4].OrderStatus)
		assert.Equal(t, "Criado", orders[5].OrderStatus)
		assert.Equal(t, "Criado", orders[6].OrderStatus)
		assert.Equal(t, "Criado", orders[7].OrderStatus)
	})

	t.Run("got success when sorting orders with few already ordered orders use case", func(t *testing.T) {
		t.Parallel()

		sut := NewSortOrdersUseCase()

		orders := []dto.OrderResponse{
			{
				OrderStatus: "Finalizado",
			},
			{
				OrderStatus: "Criado",
			},
		}

		assert.Equal(t, "Finalizado", orders[0].OrderStatus)
		assert.Equal(t, "Criado", orders[1].OrderStatus)

		sut.Execute(orders)

		assert.Equal(t, "Finalizado", orders[0].OrderStatus)
		assert.Equal(t, "Criado", orders[1].OrderStatus)
	})

	t.Run("got success when sorting orders with few unordered orders use case", func(t *testing.T) {
		t.Parallel()

		sut := NewSortOrdersUseCase()

		orders := []dto.OrderResponse{
			{
				OrderStatus: "Criado",
			},
			{
				OrderStatus: "Finalizado",
			},
		}

		assert.Equal(t, "Criado", orders[0].OrderStatus)
		assert.Equal(t, "Finalizado", orders[1].OrderStatus)

		sut.Execute(orders)

		assert.Equal(t, "Finalizado", orders[0].OrderStatus)
		assert.Equal(t, "Criado", orders[1].OrderStatus)
	})

	t.Run("got success when sorting orders with few unordered unfinished orders use case", func(t *testing.T) {
		t.Parallel()

		sut := NewSortOrdersUseCase()

		orders := []dto.OrderResponse{
			{
				OrderStatus: "Criado",
			},
			{
				OrderStatus: "Preparando",
			},
		}

		assert.Equal(t, "Criado", orders[0].OrderStatus)
		assert.Equal(t, "Preparando", orders[1].OrderStatus)

		sut.Execute(orders)

		assert.Equal(t, "Preparando", orders[0].OrderStatus)
		assert.Equal(t, "Criado", orders[1].OrderStatus)
	})
}
