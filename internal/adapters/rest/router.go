// internal/adapters/rest/router.go

package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	swagger "github.com/swaggo/fiber-swagger"

	_ "github.com/krud3/prueba-tecnica/docs"
)

func SetUpRoutes(app *fiber.App, customerHandler *CustomerHandler, workOrderHandler *WorkOrderHandler) {
	// display on console petitions using fiber logger middleware+
	app.Use(logger.New())

	// main route of api
	api := app.Group("/api/v1")

	// health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	// swagger endpoint
	api.Get("/swagger/*", swagger.WrapHandler)

	// ----- CUSTOMER
	customers := api.Group("/customers")
	customers.Post("/", customerHandler.Create)
	customers.Get("/active", customerHandler.GetActive)
	customers.Get("/all", customerHandler.GetAll)
	customers.Get("/:id", customerHandler.GetByID)

	// ----- WORKORDER
	workOrders := api.Group("/work-orders")
	workOrders.Post("/", workOrderHandler.Create)
	workOrders.Get("/", workOrderHandler.GetFiltered)
	workOrders.Get("/:id", workOrderHandler.GetByID)
	workOrders.Patch("/:id/complete", workOrderHandler.CompleteOrder)

	// ----- GET ALL ORDERS FROM A CLIENT
	customers.Get("/:customerID/work-orders", workOrderHandler.GetByCustomerID)
}
