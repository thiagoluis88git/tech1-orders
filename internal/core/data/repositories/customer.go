package repositories

import (
	"context"

	"github.com/thiagoluis88git/tech1-orders/internal/core/data/remote"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/repository"
)

type CustomerRepositoryImpl struct {
	ds remote.CustomerRemoteDataSource
}

func NewCustomerRepository(ds remote.CustomerRemoteDataSource) repository.CustomerRepository {
	return &CustomerRepositoryImpl{
		ds: ds,
	}
}

func (repo *CustomerRepositoryImpl) GetCustomerByCPF(ctx context.Context, cpf string) (dto.Customer, error) {
	response, err := repo.ds.GetCustomerByCPF(ctx, cpf)

	if err != nil {
		return dto.Customer{}, err
	}

	return dto.Customer{
		ID:    response.ID,
		Name:  response.Name,
		CPF:   response.CPF,
		Email: response.Email,
	}, nil
}
