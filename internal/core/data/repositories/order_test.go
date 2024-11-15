package repositories_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/model"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/repositories"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
)

func TestOrderRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestCreateOrderWithSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}
	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)
}

func (suite *RepositoryTestSuite) TestCreatePayingOrderWithSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}
	orderResponse, err := repo.CreatePayingOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)
}

func (suite *RepositoryTestSuite) TestDeleteOrderWithSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreatePayingOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	err = repo.DeleteOrder(suite.ctx, orderResponse.OrderId)
	suite.NoError(err)
}

func (suite *RepositoryTestSuite) TestFinishOrderWithPaymentSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	err = repo.FinishOrderWithPayment(suite.ctx, orderResponse.OrderId, uint(3))
	suite.NoError(err)
}

func (suite *RepositoryTestSuite) TestGetOrderByIDSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	cpf := "12345678910"

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		CPF:          &cpf,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	customerDS.On("GetCustomerByCPF", suite.ctx, cpf).Return(MockCustomer(), nil)

	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	orderResult, err := repo.GetOrderById(suite.ctx, orderResponse.OrderId)
	suite.NoError(err)
	suite.Equal(12, orderResult.TicketNumber)
	suite.Equal("CustomerName", *orderResult.CustomerName)
}

func (suite *RepositoryTestSuite) TestGetOrderByIDWithoutCustomerSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	orderResult, err := repo.GetOrderById(suite.ctx, orderResponse.OrderId)
	suite.NoError(err)
	suite.Equal(12, orderResult.TicketNumber)
	suite.Nil(orderResult.CustomerName)
}

func (suite *RepositoryTestSuite) TestGetOrderByIDWithUnknownIDError() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	orderResult, err := repo.GetOrderById(suite.ctx, uint(333))
	suite.Error(err)
	suite.Empty(orderResult)
}

func (suite *RepositoryTestSuite) TestGetOrdersToPrepareSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	cpf := "12345678910"

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		CPF:          &cpf,
		TicketNumber: 12,

		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	customerDS.On("GetCustomerByCPF", suite.ctx, cpf).Return(MockCustomer(), nil)

	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	ordersToPrepare, err := repo.GetOrdersToPrepare(suite.ctx)
	suite.NoError(err)
	suite.Equal(1, len(ordersToPrepare))
	suite.Equal("CustomerName", *ordersToPrepare[0].CustomerName)
}

func (suite *RepositoryTestSuite) TestGetOrdersToPrepareWithoutCustomerSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	ordersToPrepare, err := repo.GetOrdersToPrepare(suite.ctx)
	suite.NoError(err)
	suite.Equal(1, len(ordersToPrepare))
	suite.Nil(ordersToPrepare[0].CustomerName)
}

func (suite *RepositoryTestSuite) TestGetOrdersToFollowSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreateOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	err = repo.UpdateToPreparing(suite.ctx, orderResponse.OrderId)
	suite.NoError(err)

	ordersToFollow, err := repo.GetOrdersToFollow(suite.ctx)
	suite.NoError(err)
	suite.Equal(1, len(ordersToFollow))
}

func (suite *RepositoryTestSuite) TestGetOrdersWaitingPaymentSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreatePayingOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	ordersWaitingPayment, err := repo.GetOrdersWaitingPayment(suite.ctx)
	suite.NoError(err)
	suite.Equal(1, len(ordersWaitingPayment))
}

func (suite *RepositoryTestSuite) TestUpdateToDoneSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreatePayingOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	err = repo.UpdateToDone(suite.ctx, orderResponse.OrderId)
	suite.NoError(err)
}

func (suite *RepositoryTestSuite) TestUpdateToDeliveredSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreatePayingOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	err = repo.UpdateToDelivered(suite.ctx, orderResponse.OrderId)
	suite.NoError(err)
}

func (suite *RepositoryTestSuite) TestUpdateToNotDeliveredSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repoProduct := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repoProduct.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)
	newOrder := dto.Order{
		TotalPrice:   5090,
		PaymentID:    uint(12),
		TicketNumber: 12,
		OrderProduct: []dto.OrderProduct{
			{
				ProductID: uint(1),
			},
		},
	}

	orderResponse, err := repo.CreatePayingOrder(suite.ctx, newOrder)
	suite.NoError(err)
	suite.Equal(uint(1), orderResponse.OrderId)

	err = repo.UpdateToNotDelivered(suite.ctx, orderResponse.OrderId)
	suite.NoError(err)
}

func (suite *RepositoryTestSuite) TestGetNextTicketNumberSuccess() {
	// ensure that the postgres database is empty
	var tickets []model.OrderTicketNumber
	result := suite.db.Connection.Find(&tickets)
	suite.NoError(result.Error)
	suite.Empty(tickets)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)

	date := time.Now().UnixMilli()

	ticket := repo.GetNextTicketNumber(suite.ctx, date)

	suite.Equal(1, ticket)
}

func (suite *RepositoryTestSuite) TestGetNextTicketNumberWithPlus1Success() {
	// ensure that the postgres database is empty
	var tickets []model.OrderTicketNumber
	result := suite.db.Connection.Find(&tickets)
	suite.NoError(result.Error)
	suite.Empty(tickets)

	customerDS := new(MockCustomerRemoteDataSource)

	repo := repositories.NewOrderRespository(suite.db, customerDS)

	date := time.Now().UnixMilli()

	ticket := repo.GetNextTicketNumber(suite.ctx, date)

	suite.Equal(1, ticket)

	newTicket := repo.GetNextTicketNumber(suite.ctx, date)

	suite.Equal(2, newTicket)
}
