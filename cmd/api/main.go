package main

import (
	"fmt"
	"net/http"

	"github.com/thiagoluis88git/tech1-orders/internal/core/data/remote"
	"github.com/thiagoluis88git/tech1-orders/internal/core/data/repositories"
	"github.com/thiagoluis88git/tech1-orders/internal/core/domain/usecases"
	"github.com/thiagoluis88git/tech1-orders/internal/core/handler"
	"github.com/thiagoluis88git/tech1-orders/pkg/database"
	"github.com/thiagoluis88git/tech1-orders/pkg/environment"
	"github.com/thiagoluis88git/tech1-orders/pkg/httpserver"
	"github.com/thiagoluis88git/tech1-orders/pkg/responses"
	"gorm.io/driver/postgres"

	"github.com/mvrilo/go-redoc"

	"github.com/go-chi/chi/v5"

	_ "github.com/thiagoluis88git/tech1-orders/docs"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Tech1 API Docs
// @version 1.0
// @description This is the API for the Tech1 Fiap Project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localshot:3210
// @BasePath /
func main() {
	environment.LoadEnvironmentVariables()

	doc := redoc.Redoc{
		Title:       "Example API",
		Description: "Example API Description",
		SpecFile:    *environment.RedocFolderPath,
		SpecPath:    "/docs/swagger.json",
		DocsPath:    "/docs",
	}

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v",
		environment.GetDBHost(),
		environment.GetDBUser(),
		environment.GetDBPassword(),
		environment.GetDBName(),
		environment.GetDBPort(),
	)

	db, err := database.ConfigDatabase(postgres.Open(dsn))

	if err != nil {
		panic(fmt.Sprintf("could not open database: %v", err.Error()))
	}

	router := chi.NewRouter()
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Recoverer)

	httpClient := httpserver.NewHTTPClient()

	customerRemote := remote.NewCustomerRemoteDataSource(httpClient, environment.GetCustomerRootAPI())
	customerRepo := repositories.NewCustomerRepository(customerRemote)

	productRepo := repositories.NewProductRepository(db)
	validateProductCategoryUseCase := usecases.NewValidateProductCategoryUseCase()
	getCategoriesUseCase := usecases.NewGetCategoriesUseCase(productRepo)
	getProductsUseCase := usecases.NewGetProductsByCategoryUseCase(productRepo)
	getProductByIdUseCase := usecases.NewGetProductByIdUseCase(productRepo)
	deleteProductUseCase := usecases.NewDeleteProductUseCase(productRepo)
	updateProductUseCase := usecases.NewUpdateProductUseCase(productRepo)
	createProductUseCase := usecases.NewCreateProductUseCase(validateProductCategoryUseCase, productRepo)

	orderRepo := repositories.NewOrderRespository(db, customerRemote)
	validateToPreare := usecases.NewValidateOrderToPrepareUseCase(orderRepo)
	validateToDone := usecases.NewValidateOrderToDoneUseCase(orderRepo)
	validateToDeliveredOrNot := usecases.NewValidateOrderToDeliveredOrNotUseCase(orderRepo)
	sortOrders := usecases.NewSortOrdersUseCase()
	createOrderUseCase := usecases.NewCreateOrderUseCase(
		orderRepo,
		customerRepo,
		validateToPreare,
		validateToDone,
		validateToDeliveredOrNot,
		sortOrders,
	)
	getOrderByIdUseCase := usecases.NewGetOrderByIdUseCase(orderRepo)
	getOrdersToPrepareUseCase := usecases.NewGetOrdersToPrepareUseCase(
		orderRepo,
		sortOrders,
	)
	getOrdersToFollowUseCase := usecases.NewGetOrdersToFollowUseCase(
		orderRepo,
		sortOrders,
	)
	getOrdersWaitingPaymentUseCase := usecases.NewGetOrdersWaitingPaymentUseCase(
		orderRepo,
		sortOrders,
	)
	updateToPreparingUseCase := usecases.NewUpdateToPreparingUseCase(
		orderRepo,
		validateToPreare,
	)
	updateToDoneUseCase := usecases.NewUpdateToDoneUseCase(
		orderRepo,
		validateToDone,
	)
	updateToDeliveredUseCase := usecases.NewUpdateToDeliveredUseCase(
		orderRepo,
		validateToDeliveredOrNot,
	)
	updateToNotDeliveredUseCase := usecases.NewUpdateToNotDeliveredUseCase(
		orderRepo,
		validateToDeliveredOrNot,
	)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpserver.SendResponseSuccess(w, &responses.BusinessResponse{
			StatusCode: 200,
			Message:    "ok",
		})
	})

	router.Post("/api/admin/products", handler.CreateProductHandler(createProductUseCase))
	router.Delete("/api/admin/products/{id}", handler.DeleteProductHandler(deleteProductUseCase))
	router.Put("/api/admin/products/{id}", handler.UpdateProductHandler(updateProductUseCase))
	router.Get("/api/products/{id}", handler.GetProductsByIdHandler(getProductByIdUseCase))
	router.Get("/api/products/categories", handler.GetCategoriesHandler(getCategoriesUseCase))
	router.Get("/api/products/categories/{category}", handler.GetProductsByCategoryHandler(getProductsUseCase))

	router.Post("/api/orders", handler.CreateOrderHandler(createOrderUseCase))
	router.Get("/api/orders/{id}", handler.GetOrderByIdHandler(getOrderByIdUseCase))
	router.Get("/api/orders/to-prepare", handler.GetOrdersToPrepareHandler(getOrdersToPrepareUseCase))
	router.Get("/api/orders/follow", handler.GetOrdersToFollowHandler(getOrdersToFollowUseCase))
	router.Get("/api/orders/waiting-payment", handler.GetOrdersWaitingPaymentHandler(getOrdersWaitingPaymentUseCase))
	router.Put("/api/orders/{id}/preparing", handler.UpdateOrderPreparingHandler(updateToPreparingUseCase))
	router.Put("/api/orders/{id}/done", handler.UpdateOrderDoneHandler(updateToDoneUseCase))
	router.Put("/api/orders/{id}/delivered", handler.UpdateOrderDeliveredHandler(updateToDeliveredUseCase))
	router.Put("/api/orders/{id}/not-delivered", handler.UpdateOrderNotDeliveredandler(updateToNotDeliveredUseCase))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3210/swagger/doc.json"),
	))

	go http.ListenAndServe(":3211", doc.Handler())

	server := httpserver.New(router)
	server.Start()
}
