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

// Create crea un nuevo cliente.
// @Summary      Crea un nuevo cliente
// @Description  Crea un nuevo cliente en la base de datos con estado inactivo por defecto.
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer body CreateCustomerRequest true "Datos del Cliente a crear"
// @Success      201 {object} domain.Customer
// @Failure      400 {object} map[string]string "Error: Petición inválida"
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /customers [post]
func (cH *CustomerHandler) Create(c *fiber.Ctx) error {
	// using dto now
	var req CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cuerpo de la petición inválido"})
	}
	// map DTO to domain.Customer
	customer := domain.Customer{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Address:   req.Address,
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

// GetByID busca un cliente por su ID.
// @Summary      Busca un cliente por ID
// @Description  Obtiene los detalles de un cliente específico usando su UUID.
// @Tags         customers
// @Produce      json
// @Param        id path string true "ID del Cliente (UUID)"
// @Success      200 {object} domain.Customer
// @Failure      400 {object} map[string]string "Error: ID inválido"
// @Failure      404 {object} map[string]string "Error: Cliente no encontrado"
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /customers/{id} [get]
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

// GetActive obtiene todos los clientes activos.
// @Summary      Obtiene clientes activos
// @Description  Devuelve una lista de todos los clientes cuyo estado es 'is_active = true'.
// @Tags         customers
// @Produce      json
// @Success      200 {array} domain.Customer
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /customers/active [get]
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
