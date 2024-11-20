package repositories_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/model"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/repositories"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func TestCustomerRepository(t *testing.T) {
	t.Parallel()

	t.Run("got success when calling GetCustomerByCPF repository", func(t *testing.T) {

		ds := new(MockCustomerRemoteDataSource)
		sut := repositories.NewCustomerRepository(ds)

		ctx := context.TODO()

		ds.On("GetCustomerByCPF", ctx, "12345678910").Return(model.Customer{
			ID:   uint(2),
			Name: "Name",
		}, nil)

		response, err := sut.GetCustomerByCPF(ctx, "12345678910")

		assert.NoError(t, err)
		assert.NotEmpty(t, response)
		assert.Equal(t, "Name", response.Name)
	})

	t.Run("got error when calling GetCustomerByCPF repository", func(t *testing.T) {

		ds := new(MockCustomerRemoteDataSource)
		sut := repositories.NewCustomerRepository(ds)

		ctx := context.TODO()

		ds.On("GetCustomerByCPF", ctx, "12345678910").Return(model.Customer{}, &responses.NetworkError{
			Code: 500,
		})

		response, err := sut.GetCustomerByCPF(ctx, "12345678910")

		assert.Error(t, err)
		assert.Empty(t, response)
	})
}
