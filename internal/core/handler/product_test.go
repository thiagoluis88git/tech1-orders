package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/dto"
	"github.com/thiagoluis88git/tech1-orders/internal/core/handler"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
)

func mockProduct() dto.ProductForm {
	return dto.ProductForm{
		Name:        "Name",
		Description: "Description",
		Category:    "Category",
		Price:       12.55,
		Images: []dto.ProducImage{
			{
				ImageUrl: "ImageUrl",
			},
		},
	}
}

func mockUpdateProduct() dto.ProductForm {
	return dto.ProductForm{
		Id:          uint(12),
		Name:        "Name",
		Description: "Description",
		Category:    "Category",
		Price:       12.55,
		Images: []dto.ProducImage{
			{
				ImageUrl: "ImageUrl",
			},
		},
	}
}

func TestCreateProductHandler(t *testing.T) {
	t.Parallel()

	t.Run("got success when calling get category handler handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockProduct())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/products/categories", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getCategoryUseCase := new(MockGetCategoryUseCase)

		getCategoryUseCase.On("Execute").Return([]string{"CAT1, CAT2"})

		getCategoryHandler := handler.GetCategoriesHandler(getCategoryUseCase)

		getCategoryHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response []string
		err = json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NoError(t, err)

		assert.Equal(t, "CAT1, CAT2", response[0])
	})

	t.Run("got success when calling create product handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockProduct())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/products", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		createProductUseCase := new(MockCreateProductUseCase)

		createProductUseCase.On("Execute", req.Context(), mockProduct()).
			Return(uint(2), nil)

		loginCustomerHandler := handler.CreateProductHandler(createProductUseCase)

		loginCustomerHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.ProductCreationResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NoError(t, err)

		assert.Equal(t, uint(2), response.Id)
	})

	t.Run("got error on CreateProduct UseCase when calling create product handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockProduct())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/products", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		createProductUseCase := new(MockCreateProductUseCase)

		createProductUseCase.On("Execute", req.Context(), mockProduct()).
			Return(uint(0), &responses.BusinessResponse{
				StatusCode: 422,
			})

		loginCustomerHandler := handler.CreateProductHandler(createProductUseCase)

		loginCustomerHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})

	t.Run("got error on invalid json UseCase when calling create product handler", func(t *testing.T) {
		t.Parallel()

		body := bytes.NewBuffer([]byte("asff{{}"))

		req := httptest.NewRequest(http.MethodPost, "/api/products", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		createProductUseCase := new(MockCreateProductUseCase)

		createProductUseCase.On("Execute", req.Context(), mockProduct()).
			Return(uint(0), &responses.BusinessResponse{
				StatusCode: 422,
			})

		loginCustomerHandler := handler.CreateProductHandler(createProductUseCase)

		loginCustomerHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got success when calling get products by category handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/products/category/{category}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("category", "CATEGORY")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getProductsByCategoryUseCase := new(MockGetProductsByCategoryUseCase)

		getProductsByCategoryUseCase.On("Execute", req.Context(), "CATEGORY").
			Return([]dto.ProductResponse{
				{
					Id:   uint(1),
					Name: "Name 1",
				},
				{
					Id:   uint(2),
					Name: "Name 2",
				},
			}, nil)

		getProductsByCategoryHandler := handler.GetProductsByCategoryHandler(getProductsByCategoryUseCase)

		getProductsByCategoryHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response []dto.ProductResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NoError(t, err)

		assert.Equal(t, uint(1), response[0].Id)
		assert.Equal(t, uint(2), response[1].Id)
		assert.Equal(t, "Name 1", response[0].Name)
		assert.Equal(t, "Name 2", response[1].Name)
	})

	t.Run("got error on GetProducts UseCase when calling get products by category handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/products/category/{category}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("category", "CATEGORY")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getProductsByCategoryUseCase := new(MockGetProductsByCategoryUseCase)

		getProductsByCategoryUseCase.On("Execute", req.Context(), "CATEGORY").
			Return([]dto.ProductResponse{}, &responses.BusinessResponse{
				StatusCode: 422,
			})

		getProductsByCategoryHandler := handler.GetProductsByCategoryHandler(getProductsByCategoryUseCase)

		getProductsByCategoryHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})

	t.Run("got error on missing category when calling get products by category handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/products/category/{category}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getProductsByCategoryUseCase := new(MockGetProductsByCategoryUseCase)

		getProductsByCategoryUseCase.On("Execute", req.Context(), "CATEGORY").
			Return([]dto.ProductResponse{}, &responses.BusinessResponse{
				StatusCode: 422,
			})

		getProductsByCategoryHandler := handler.GetProductsByCategoryHandler(getProductsByCategoryUseCase)

		getProductsByCategoryHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got success when calling get products by id handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "3")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getProductsByIDUseCase := new(MockGetProductsByIDUseCase)

		getProductsByIDUseCase.On("Execute", req.Context(), uint(3)).
			Return(dto.ProductResponse{
				Id:   uint(1),
				Name: "Name 1",
			}, nil)

		getProductsByIDHandler := handler.GetProductsByIdHandler(getProductsByIDUseCase)

		getProductsByIDHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response dto.ProductResponse
		err := json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.NoError(t, err)

		assert.Equal(t, uint(1), response.Id)
		assert.Equal(t, "Name 1", response.Name)
	})

	t.Run("got error on GetProducts UseCase when calling get products by id handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "3")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getProductsByIDUseCase := new(MockGetProductsByIDUseCase)

		getProductsByIDUseCase.On("Execute", req.Context(), uint(3)).
			Return(dto.ProductResponse{}, &responses.BusinessResponse{
				StatusCode: 422,
			}).Times(1)

		getProductsByIDHandler := handler.GetProductsByIdHandler(getProductsByIDUseCase)

		getProductsByIDHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})

	t.Run("got error on invalid id when calling get products by id handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "x3")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getProductsByIDUseCase := new(MockGetProductsByIDUseCase)

		getProductsByIDHandler := handler.GetProductsByIdHandler(getProductsByIDUseCase)

		getProductsByIDHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got error on missing id when calling get products by id handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		getProductsByIDUseCase := new(MockGetProductsByIDUseCase)

		getProductsByIDHandler := handler.GetProductsByIdHandler(getProductsByIDUseCase)

		getProductsByIDHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got success when calling deleting product handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodDelete, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "3")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		deleteProductUseCase := new(MockDeleteProductUseCase)

		deleteProductUseCase.On("Execute", req.Context(), uint(3)).
			Return(nil)

		deleteProductHandler := handler.DeleteProductHandler(deleteProductUseCase)

		deleteProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("got error on DeleteProduct UseCase when calling deleting product handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodDelete, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "3")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		deleteProductUseCase := new(MockDeleteProductUseCase)

		deleteProductUseCase.On("Execute", req.Context(), uint(3)).
			Return(&responses.BusinessResponse{
				StatusCode: 422,
			})

		deleteProductHandler := handler.DeleteProductHandler(deleteProductUseCase)

		deleteProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})

	t.Run("got error on invalid id when calling deleting product handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodDelete, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "s3")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		deleteProductUseCase := new(MockDeleteProductUseCase)

		deleteProductUseCase.On("Execute", req.Context(), uint(3)).
			Return(&responses.BusinessResponse{
				StatusCode: 422,
			})

		deleteProductHandler := handler.DeleteProductHandler(deleteProductUseCase)

		deleteProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got error on missing id when calling deleting product handler", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodDelete, "/api/products/{id}", nil)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		deleteProductUseCase := new(MockDeleteProductUseCase)

		deleteProductUseCase.On("Execute", req.Context(), uint(3)).
			Return(&responses.BusinessResponse{
				StatusCode: 422,
			})

		deleteProductHandler := handler.DeleteProductHandler(deleteProductUseCase)

		deleteProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got success when calling updating product handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockUpdateProduct())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/products/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateProductUseCase := new(MockUpdateProductUseCase)

		updateProductUseCase.On("Execute", req.Context(), mockUpdateProduct()).
			Return(nil)

		updateProductHandler := handler.UpdateProductHandler(updateProductUseCase)

		updateProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run("got error on UpdateProduct UseCase when calling updating product handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockUpdateProduct())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/products/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateProductUseCase := new(MockUpdateProductUseCase)

		updateProductUseCase.On("Execute", req.Context(), mockUpdateProduct()).
			Return(&responses.BusinessResponse{
				StatusCode: 422,
			})

		updateProductHandler := handler.UpdateProductHandler(updateProductUseCase)

		updateProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})

	t.Run("got error on invalid json when calling updating product handler", func(t *testing.T) {
		t.Parallel()

		body := bytes.NewBuffer([]byte("dsgfg{"))

		req := httptest.NewRequest(http.MethodPost, "/api/products/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateProductUseCase := new(MockUpdateProductUseCase)

		updateProductHandler := handler.UpdateProductHandler(updateProductUseCase)

		updateProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got error on invalid id when calling updating product handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockUpdateProduct())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/products/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "x12")

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateProductUseCase := new(MockUpdateProductUseCase)

		updateProductHandler := handler.UpdateProductHandler(updateProductUseCase)

		updateProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("got error on missing id when calling updating product handler", func(t *testing.T) {
		t.Parallel()

		jsonData, err := json.Marshal(mockUpdateProduct())

		assert.NoError(t, err)

		body := bytes.NewBuffer(jsonData)

		req := httptest.NewRequest(http.MethodPost, "/api/products/{id}", body)
		req.Header.Add("Content-Type", "application/json")

		rctx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		recorder := httptest.NewRecorder()

		updateProductUseCase := new(MockUpdateProductUseCase)

		updateProductHandler := handler.UpdateProductHandler(updateProductUseCase)

		updateProductHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}
