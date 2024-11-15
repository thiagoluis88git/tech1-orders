package remote

import (
	"context"
	"fmt"
	"net/http"

	"github.com/thiagoluis88git/tech1-orders/internal/core/data/model"
	"github.com/thiagoluis88git/tech1-orders/pkg/httpserver"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

type CustomerRemoteDataSource interface {
	GetCustomerByCPF(ctx context.Context, cpf string) (model.Customer, error)
}

type CustomerRemoteDataSourceImpl struct {
	client  *http.Client
	rootURL string
}

func NewCustomerRemoteDataSource(client *http.Client, rootURL string) CustomerRemoteDataSource {
	return &CustomerRemoteDataSourceImpl{
		client:  client,
		rootURL: rootURL,
	}
}

func (ds *CustomerRemoteDataSourceImpl) GetCustomerByCPF(ctx context.Context, cpf string) (model.Customer, error) {
	endpoint := fmt.Sprintf("%v/%v/%v", ds.rootURL, "customer", cpf)

	response, err := httpserver.DoGetRequest(
		ctx,
		ds.client,
		endpoint,
		nil,
		model.Customer{},
	)

	if err != nil {
		return model.Customer{}, &responses.NetworkError{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	return response, nil
}
