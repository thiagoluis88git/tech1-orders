package repositories

import (
	"context"

	"github.com/thiagoluis88git/tech1-orders/internal/core/data/model"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/repository"
	"github.com/thiagoluis88git/tech1-orders/pkg/database"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *database.Database
}

func NewProductRepository(db *database.Database) repository.ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (repository *ProductRepository) GetCategories() []string {
	return []string{
		model.CategoryCombo,
		model.CategorySnack,
		model.CategoryBeverage,
		model.CategoryToppings,
		model.CategoryDesert,
	}
}

func (repository *ProductRepository) CreateProduct(ctx context.Context, product dto.ProductForm) (uint, error) {
	tx := repository.db.Connection.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return 0, responses.GetDatabaseError(err)
	}

	productEntity := &model.Product{
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Price:       product.Price,
	}

	err := tx.Create(productEntity).Error

	if err != nil {
		tx.Rollback()
		return 0, responses.GetDatabaseError(err)
	}

	productImages := []*model.ProductImage{}

	for _, value := range product.Images {
		productImages = append(productImages, &model.ProductImage{
			ProductID: productEntity.ID,
			ImageUrl:  value.ImageUrl,
		})
	}

	err = tx.Create(productImages).Error

	if err != nil {
		tx.Rollback()
		return 0, responses.GetDatabaseError(err)
	}

	err = repository.createComboIfProductsNedded(tx, product, productEntity.ID)

	if err != nil {
		tx.Rollback()
		return 0, responses.GetDatabaseError(err)
	}

	err = tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return 0, responses.GetDatabaseError(err)
	}

	return productEntity.ID, nil
}

func (repository *ProductRepository) createComboIfProductsNedded(
	tx *gorm.DB,
	productWithCombo dto.ProductForm,
	comboId uint,
) error {
	if productWithCombo.ComboProductsIds != nil {
		for _, value := range *productWithCombo.ComboProductsIds {
			comboProductEntity := &model.ComboProduct{
				ProductID:      comboId,
				ComboProductID: value,
			}

			err := tx.Create(comboProductEntity).Error

			if err != nil {
				tx.Rollback()
				return responses.GetDatabaseError(err)
			}
		}
	}

	return nil
}

func (repository *ProductRepository) GetProductsByCategory(ctx context.Context, category string) ([]dto.ProductResponse, error) {
	var productmodel []model.Product
	err := repository.
		db.Connection.WithContext(ctx).
		Model(&model.Product{}).
		Preload("ProductImage").
		Preload("ComboProduct").
		Where("category = ?", category).
		Find(&productmodel).
		Error

	if err != nil {
		return []dto.ProductResponse{}, responses.GetDatabaseError(err)
	}

	return repository.buildProducts(ctx, productmodel), nil
}

func (repository *ProductRepository) GetProductById(ctx context.Context, id uint) (dto.ProductResponse, error) {
	var productEntity model.Product
	err := repository.
		db.Connection.WithContext(ctx).
		Model(&model.Product{}).
		Preload("ProductImage").
		Preload("ComboProduct").
		First(&productEntity, id).
		Error

	if err != nil {
		return dto.ProductResponse{}, responses.GetDatabaseError(err)
	}

	return repository.buildProduct(ctx, productEntity), nil
}

func (repository *ProductRepository) DeleteProduct(ctx context.Context, productId uint) error {
	tx := repository.db.Connection.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return responses.GetDatabaseError(err)
	}

	err := tx.Where("product_id = ?", productId).Unscoped().Delete(&model.ProductImage{}).Error

	if err != nil {
		tx.Rollback()
		return responses.GetDatabaseError(err)
	}

	err = tx.Where("product_id = ?", productId).Unscoped().Delete(&model.ComboProduct{}).Error

	if err != nil {
		tx.Rollback()
		return responses.GetDatabaseError(err)
	}

	err = tx.Unscoped().Delete(&model.Product{}, productId).Error

	if err != nil {
		tx.Rollback()
		return responses.GetDatabaseError(err)
	}

	err = tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *ProductRepository) UpdateProduct(ctx context.Context, product dto.ProductForm) error {
	productEntity := model.Product{
		Model:       gorm.Model{ID: product.Id},
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Price:       product.Price,
	}

	err := repository.db.Connection.WithContext(ctx).Save(&productEntity).Error

	if err != nil {
		return responses.GetDatabaseError(err)
	}

	return nil
}

func (repository *ProductRepository) buildProducts(ctx context.Context, productmodel []model.Product) []dto.ProductResponse {
	products := []dto.ProductResponse{}

	for _, value := range productmodel {
		products = append(products, repository.buildProduct(ctx, value))
	}

	return products
}

func (repository *ProductRepository) buildProduct(ctx context.Context, value model.Product) dto.ProductResponse {
	images := []dto.ProducImage{}

	for _, valueImage := range value.ProductImage {
		images = append(images, dto.ProducImage{
			ImageUrl: valueImage.ImageUrl,
		})
	}

	comboProducts := repository.getComboProductsIfNedded(ctx, value)

	return dto.ProductResponse{
		Id:            value.ID,
		Name:          value.Name,
		Description:   value.Description,
		Category:      value.Category,
		Price:         value.Price,
		Images:        images,
		ComboProducts: comboProducts,
	}
}

func (repository *ProductRepository) getComboProductsIfNedded(ctx context.Context, value model.Product) *[]dto.ProductResponse {
	var comboProducts []dto.ProductResponse

	if value.ComboProduct != nil {
		comboProducts = make([]dto.ProductResponse, 0)

		for _, comboProduct := range value.ComboProduct {
			var product model.Product

			err := repository.db.Connection.
				WithContext(ctx).
				Preload("ProductImage").
				First(&product, comboProduct.ComboProductID).
				Error

			if err == nil {
				comboProducts = append(comboProducts, repository.buildProduct(ctx, product))
			}
		}
	}

	return &comboProducts
}
