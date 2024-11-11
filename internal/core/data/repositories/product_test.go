package repositories_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/model"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/repositories"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func TestProductRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestGetProductsWithSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repo := repositories.NewProductRepository(suite.db)

	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Lanches",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}
	newId, err := repo.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	createdProducts, err := repo.GetProductsByCategory(suite.ctx, "Lanches")

	suite.NoError(err)
	suite.Equal(1, len(createdProducts))
	suite.Equal("New Product Created", createdProducts[0].Name)
}

func (suite *RepositoryTestSuite) TestGetProductByIdWithSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repo := repositories.NewProductRepository(suite.db)

	newProduct := dto.ProductForm{
		Name:        "New Product Created",
		Description: "New Description Product Created",
		Category:    "Lanches",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}
	newId, err := repo.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	createdProduct, err := repo.GetProductById(suite.ctx, uint(1))

	suite.NoError(err)
	suite.Equal("New Product Created", createdProduct.Name)
}

func (suite *RepositoryTestSuite) TestCreateProductWithSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	repo := repositories.NewProductRepository(suite.db)
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
	newId, err := repo.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	// ensure that we have a new product in the database
	result = suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Equal(1, len(products))
	suite.Equal(uint(1), products[0].ID)
	suite.Equal("New Product Created", products[0].Name)
	suite.Equal("New Description Product Created", products[0].Description)
}

func (suite *RepositoryTestSuite) TestUpdateProductWithSuccess() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	// create repository and save new note
	repo := repositories.NewProductRepository(suite.db)
	newProduct := dto.ProductForm{
		Name:        "New Product",
		Description: "New Description Product",
		Category:    "Category",
		Price:       2990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	newId, err := repo.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	// ensure that we have a new product in the database
	result = suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Equal(1, len(products))
	suite.Equal(uint(1), products[0].ID)
	suite.Equal("New Product", products[0].Name)
	suite.Equal("New Description Product", products[0].Description)

	updateProduct := dto.ProductForm{
		Id:          uint(1),
		Name:        "Updated Product",
		Description: "Updated Description Product",
		Category:    "Category",
		Price:       3990,
		Images: []dto.ProducImage{
			{
				ImageUrl: "NewImageUrl",
			},
		},
	}

	err = repo.UpdateProduct(suite.ctx, updateProduct)

	suite.NoError(err)

	// ensure that we have a new product in the database
	result = suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Equal(1, len(products))
	suite.Equal(uint(1), products[0].ID)
	suite.Equal("Updated Product", products[0].Name)
	suite.Equal("Updated Description Product", products[0].Description)
}

func (suite *RepositoryTestSuite) TestCreateProductWithConflictError() {
	// ensure that the postgres database is empty
	var products []model.Product
	result := suite.db.Connection.Find(&products)
	suite.NoError(result.Error)
	suite.Empty(products)

	// create repository and save new note
	repo := repositories.NewProductRepository(suite.db)
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
	newId, err := repo.CreateProduct(suite.ctx, newProduct)
	suite.NoError(err)
	suite.Equal(uint(1), newId)

	newId, err = repo.CreateProduct(suite.ctx, newProduct)

	suite.Error(err)
	suite.Equal(uint(0), newId)

	var businessError *responses.LocalError
	suite.Equal(true, errors.As(err, &businessError))
	suite.Equal(responses.DATABASE_CONFLICT_ERROR, businessError.Code)
}
