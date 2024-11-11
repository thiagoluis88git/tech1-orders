package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/usecases"
	"github.com/thiagoluis88git/tech1-orders/pkg/httpserver"
)

// @Summary Create new product
// @Description Create new product
// @Tags Product
// @Accept json
// @Produce json
// @Param product body dto.ProductForm true "product"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 "Product has required fields"
// @Failure 409 "This Product is already added"
// @Router /api/admin/products [post]
func CreateProductHandler(createUseCase usecases.CreateProductUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product dto.ProductForm

		err := httpserver.DecodeJSONBody(w, r, &product)

		if err != nil {
			log.Print("decoding product body", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		productId, err := createUseCase.Execute(r.Context(), product)

		if err != nil {
			log.Print("create product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, dto.ProductCreationResponse{
			Id: productId,
		})
	}
}

// @Summary List all products by a category
// @Description List all products by a category
// @Tags Product
// @Param category path string true "Lanches"
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ProductResponse
// @Router /api/products/categories/{category} [get]
func GetProductsByCategoryHandler(getProductsUseCase usecases.GetProductsByCategoryUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, err := httpserver.GetPathParamFromRequest(r, "category")

		if err != nil {
			log.Print("get products by category", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		products, err := getProductsUseCase.Execute(r.Context(), category)

		if err != nil {
			log.Print("get products by category", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, products)
	}
}

// @Summary Get product by ID
// @Description Get product by ID
// @Tags Product
// @Param id path int true "12"
// @Accept json
// @Produce json
// @Success 200 {object} dto.ProductResponse
// @Router /api/products/{id} [get]
func GetProductsByIdHandler(getProductById usecases.GetProductByIdUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productIdStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("get product by id", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		productId, err := strconv.Atoi(productIdStr)

		if err != nil {
			log.Print("get product by id", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		product, err := getProductById.Execute(r.Context(), uint(productId))

		if err != nil {
			log.Print("get product by id", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseSuccess(w, product)
	}
}

// @Summary Delete a product
// @Description Delete a product by ID
// @Tags Product
// @Param id path int true "12"
// @Accept json
// @Produce json
// @Success 204
// @Router /api/admin/products/{id} [delete]
func DeleteProductHandler(deleteProduct usecases.DeleteProductUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productIdStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("delete product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		productId, err := strconv.Atoi(productIdStr)

		if err != nil {
			log.Print("delete product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		err = deleteProduct.Execute(r.Context(), uint(productId))

		if err != nil {
			log.Print("delete product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}

// @Summary Update a product
// @Description Update a product by ID
// @Tags Product
// @Param id path int true "12"
// @Accept json
// @Produce json
// @Success 204
// @Router /api/admin/products/{id} [put]
func UpdateProductHandler(updateProduct usecases.UpdateProductUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productIdStr, err := httpserver.GetPathParamFromRequest(r, "id")

		if err != nil {
			log.Print("update product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		productId, err := strconv.Atoi(productIdStr)

		if err != nil {
			log.Print("update product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendBadRequestError(w, err)
			return
		}

		var product dto.ProductForm

		err = httpserver.DecodeJSONBody(w, r, &product)

		if err != nil {
			log.Print("decoding product body for update product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		product.Id = uint(productId)
		err = updateProduct.Execute(r.Context(), product)

		if err != nil {
			log.Print("update product", map[string]interface{}{
				"error":  err.Error(),
				"status": httpserver.GetStatusCodeFromError(err),
			})
			httpserver.SendResponseError(w, err)
			return
		}

		httpserver.SendResponseNoContentSuccess(w)
	}
}

// @Summary Get all categories
// @Description Get all categories to filter in products by category
// @Tags Product
// @Param id path int true "12"
// @Accept json
// @Produce json
// @Success 200 {object} []string
// @Router /api/products/categories [get]
func GetCategoriesHandler(getCategoriesUseCase usecases.GetCategoriesUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpserver.SendResponseSuccess(w, getCategoriesUseCase.Execute())
	}
}
