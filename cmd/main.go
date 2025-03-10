package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"order_management/internal/handlers"
	"order_management/internal/repositories"
	"order_management/internal/services"
	"order_management/internal/validators"
	"order_management/pkg/database"
)

func main() {
	// Initialize database
	db := database.InitDB()
	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Initialize repositories
	productRepo := repositories.NewProductRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Initialize services
	productService := services.NewProductService(productRepo, db)
	orderService := services.NewOrderService(orderRepo, productRepo, db)

	// Initialize Echo and middleware
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configurar el validador globalmente
	e.Validator = validators.NewValidator()

	// Define el grupo con el prefijo "/api"
	apiGroup := e.Group("/api")

	// Register handlers
	handlers.NewProductHandler(apiGroup, productService)
	handlers.NewOrderHandler(apiGroup, orderService, redisClient)

	e.Logger.Fatal(e.Start(":8080"))
}
