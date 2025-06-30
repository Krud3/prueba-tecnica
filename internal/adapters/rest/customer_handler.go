// internal/adapters/rest/customer_handler.go

package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/krud3/prueba-tecnica/internal/core/domain"
	"github.com/krud3/prueba-tecnica/internal/core/services"
)

type CustomerHandler struct {
	cS *services.CustomerService
}

// builder
func NewCustomerHandler(cS *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{cS: cS}
}

// POST /customers
func (cH *CustomerHandler) Create(c *fiber.Ctx) error {
	var customer domain.Customer
	// check if params match domain.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cuerpo de la petición inválido"})
	}
	// try to create
	err := cH.cS.Create(c.Context(), customer)
	if err != nil {
		// 500 server error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// 201 created
	return c.Status(fiber.StatusCreated).JSON(customer)
}

// GET /customers/:id
func (cH *CustomerHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	customerID, err := uuid.Parse(idStr)
	if err != nil {
		// err ID
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "El campo ID de la orden es inválido"})
	}
	// handler to get service to find by id
	customer, err := cH.cS.FindByID(c.Context(), customerID)
	if err != nil {
		// 500 server error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if customer == nil {
		// 404 not found
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "orden no encontrada"})
	}
	// 200 ok
	return c.Status(fiber.StatusOK).JSON(customer)
}

// GET /customers/active
func (cH *CustomerHandler) GetActive(c *fiber.Ctx) error {
	// using handler to get the service to get actives
	customers, err := cH.cS.GetActive(c.Context())
	if err != nil {
		// 500 server error due user can not send invalid data
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error al buscar los clientes activos"})
	}
	// 200 ok or empty
	return c.Status(fiber.StatusOK).JSON(customers)
}
