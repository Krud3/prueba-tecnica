// cmd/api/main.go
package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/krud3/prueba-tecnica/internal/adapters/rest"
	"github.com/krud3/prueba-tecnica/internal/adapters/storage"
	"github.com/krud3/prueba-tecnica/internal/core/services"
)

// @title API de Órdenes de Servicio
// @version 1.0
// @description Esta es la API para la prueba técnica de Fullstack.
// @termsOfService http://swagger.io/terms/
// @host localhost:3000
// @BasePath /api/v1
func main() {

	// get env
	if err := godotenv.Load(); err != nil {
		log.Println("Por favor suministrar el .env en el root de la manera en que .env.example lo dice.")
	}

	db, err := storage.NewGormDB()
	if err != nil {
		// like printf but ends with exit(0)
		log.Fatalf("Error conectando la base de datos: %v", err)
	}
	log.Println("Conexión establecida con la base de datos.")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Println("Por favor suministrar el .env en el root de la manera en que .env.example lo dice.")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	// test conection with redis
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Error contectando a Redis: %v", err)
	}
	log.Println("Conectado a Redis.")

	// create repository
	customerRepo := storage.NewGormCustomerRepository(db)
	workOrderRepo := storage.NewGormWorkOrderRepository(db)

	// stream for redis
	streamName := "work_orders_stream"
	// create services passing repositories
	customerService := services.NewCustomerService(customerRepo)
	workOrderService := services.NewWorkOrderService(workOrderRepo, customerRepo, redisClient, streamName)

	// create API handlers passing services
	customerHandler := rest.NewCustomerHandler(customerService)
	workOrderHandler := rest.NewWorkOrderHandler(workOrderService)

	// create web server with fiber
	app := fiber.New()

	// allows vite to make petitions
	allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigin,
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// config routes from API, calls handlers
	rest.SetUpRoutes(app, customerHandler, workOrderHandler)

	// init server
	port := "3000"
	log.Printf("Servidor escuchando en el puerto :%s", port)
	log.Fatal(app.Listen(":" + port))

}
