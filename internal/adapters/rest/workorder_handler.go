// internal/adapters/rest/workorder_handler.go

package rest

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/krud3/prueba-tecnica/internal/core/domain"
	"github.com/krud3/prueba-tecnica/internal/core/ports"
	"github.com/krud3/prueba-tecnica/internal/core/services"
	"gorm.io/gorm"
)

type WorkOrderHandler struct {
	wS *services.WorkOrderService
}

// builder
func NewWorkOrderHandler(wS *services.WorkOrderService) *WorkOrderHandler {
	return &WorkOrderHandler{wS: wS}
}

// Create crea una nueva orden de trabajo.
// @Summary      Crea una nueva orden de trabajo
// @Description  Crea una nueva orden para un cliente. Valida reglas de negocio como el estado del cliente y el intervalo de fechas.
// @Tags         work-orders
// @Accept       json
// @Produce      json
// @Param        workOrder body CreateWorkOrderRequest true "Datos de la Orden de Trabajo a crear"
// @Success      201 {object} domain.WorkOrder
// @Failure      400 {object} map[string]string "Error: Petición inválida"
// @Failure      404 {object} map[string]string "Error: Cliente no encontrado"
// @Failure      409 {object} map[string]string "Error: Conflicto de negocio (ej. cliente ya activo)"
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /work-orders [post]
func (wH *WorkOrderHandler) Create(c *fiber.Ctx) error {
	// now with DTO
	var req CreateWorkOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cuerpo de la petición inválido"})
	}

	workOrder := domain.WorkOrder{
		CustomerID:       req.CustomerID,
		Description:      req.Description,
		PlannedDateBegin: req.PlannedDateBegin,
		PlannedDateEnd:   req.PlannedDateEnd,
		Type:             req.Type,
	}
	// using handler to get the service to create workOrder
	err := wH.wS.Create(c.Context(), workOrder)
	if err != nil {
		switch {
		// handle custom errors
		case errors.Is(err, services.ErrAA), errors.Is(err, services.ErrCC), errors.Is(err, services.ErrDateIntertal):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		// handle customer not found
		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "el cliente especificado no existe"})
		default:
			// 500 default error
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	// 201 created
	return c.Status(fiber.StatusCreated).JSON(workOrder)
}

// CompleteOrder completa una orden de trabajo.
// @Summary      Completa una orden de trabajo
// @Description  Marca una orden como 'done', lo que activa/desactiva al cliente asociado y envía un evento a Redis.
// @Tags         work-orders
// @Produce      json
// @Param        id path string true "ID de la Orden de Trabajo (UUID)"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string "Error: ID inválido"
// @Failure      404 {object} map[string]string "Error: Orden no encontrada"
// @Failure      409 {object} map[string]string "Error: Conflicto de estado (ej. la orden ya está completada)"
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /work-orders/{id}/complete [patch]
func (wH *WorkOrderHandler) CompleteOrder(c *fiber.Ctx) error {
	idStr := c.Params("id")
	workOrderID, err := uuid.Parse(idStr)
	// verifies if id match uuid struct
	if err != nil {
		// 400
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "EL campo ID de la orden es inválido"})
	}

	// the service try to CompleteOrder
	err = wH.wS.CompleteOrder(c.Context(), workOrderID)
	if err != nil {
		switch {
		// custom errors
		case errors.Is(err, services.ErrWODone), errors.Is(err, services.ErrWOCancelled):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		// not found
		case errors.Is(err, gorm.ErrRecordNotFound):
			// 404
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "orden no encontrada"})
		default:
			// 500
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	// 200 ok
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Orden completada exitosamente"})
}

// GetByID busca una orden de trabajo por su ID.
// @Summary      Busca una orden de trabajo por ID
// @Description  Obtiene los detalles de una orden de trabajo, incluyendo la información del cliente embebida.
// @Tags         work-orders
// @Produce      json
// @Param        id path string true "ID de la Orden de Trabajo (UUID)"
// @Success      200 {object} domain.WorkOrder
// @Failure      400 {object} map[string]string "Error: ID inválido"
// @Failure      404 {object} map[string]string "Error: Orden no encontrada"
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /work-orders/{id} [get]
func (wH *WorkOrderHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	workOrderID, err := uuid.Parse(idStr)
	// verifies if id match uuid struct
	if err != nil {
		// 400
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "EL campo ID de la orden es inválido"})
	}

	workOrder, err := wH.wS.FindByID(c.Context(), workOrderID)
	// handle error
	if err != nil {
		//500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// empty?
	if workOrder == nil {
		// 404
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "orden no encontrada"})
	}

	// 200 ok
	return c.Status(fiber.StatusOK).JSON(workOrder)
}

// GetFiltered busca órdenes de trabajo con filtros.
// @Summary      Busca órdenes de trabajo con filtros
// @Description  Obtiene una lista de órdenes de trabajo. Se puede filtrar por rango de fechas (since, until) y/o por estado (status).
// @Tags         work-orders
// @Produce      json
// @Param        since  query string false "Fecha de inicio (Formato RFC3339: 2024-07-30T10:00:00Z)"
// @Param        until  query string false "Fecha de fin (Formato RFC3339: 2024-07-30T10:00:00Z)"
// @Param        status query string false "Estado de la orden" Enums(new, done, cancelled)
// @Success      200 {array} domain.WorkOrder
// @Failure      400 {object} map[string]string "Error: Parámetro de filtro inválido"
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /work-orders [get]
func (wH *WorkOrderHandler) GetFiltered(c *fiber.Ctx) error {
	// struct ports.WorkOrderFilters
	filters := ports.WorkOrderFilters{}
	// get since value
	sinceStr := c.Query("since")
	if sinceStr != "" {
		// verify time format
		since, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			// 400
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "formato de fecha 'since' inválido, usar formato RFC3339 (YYYY-MM-DDTHH:MM:SSZ)",
			})
		}
		filters.Since = &since
	}

	// get until value
	untilStr := c.Query("until")
	if untilStr != "" {
		// verify time format
		until, err := time.Parse(time.RFC3339, untilStr)
		if err != nil {
			// 400
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "formato de fecha 'until' inválido, usar formato RFC3339 (YYYY-MM-DDTHH:MM:SSZ)",
			})
		}
		filters.Until = &until
	}

	// verifies if since > until
	if filters.Since != nil && filters.Until != nil && filters.Since.After(*filters.Until) {
		// 400
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "since no puede ser posterior a until"})
	}
	// get status value
	statusStr := c.Query("status")
	if statusStr != "" {
		status := domain.Status(statusStr)
		// verifies if status is valid
		if status != domain.StatusNew && status != domain.StatusDone && status != domain.StatusCancelled {
			// 400
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "valor de 'status' inválido, debe ser 'new', 'done' o 'cancelled'",
			})
		}
		filters.Status = &status
	}

	// trying to find by filter using service
	workOrders, err := wH.wS.FindByFilter(c.Context(), filters)
	if err != nil {
		// 500 server error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error al buscar las órdenes de trabajo"})
	}

	// 200 ok
	return c.Status(fiber.StatusOK).JSON(workOrders)
}

// GetByCustomerID busca órdenes de trabajo por ID de cliente.
// @Summary      Busca órdenes de trabajo por ID de cliente
// @Description  Obtiene una lista de todas las órdenes de trabajo asociadas a un cliente específico.
// @Tags         work-orders, customers
// @Produce      json
// @Param        customerID path string true "ID del Cliente (UUID)"
// @Success      200 {array} domain.WorkOrder
// @Failure      400 {object} map[string]string "Error: ID de cliente inválido"
// @Failure      500 {object} map[string]string "Error: Error interno del servidor"
// @Router       /customers/{customerID}/work-orders [get]
func (wH *WorkOrderHandler) GetByCustomerID(c *fiber.Ctx) error {
	idStr := c.Params("customerID")
	customerID, err := uuid.Parse(idStr)
	if err != nil {
		// id given must match uuid struct
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID de cliente inválido"})
	}

	// trying to find using service
	workOrders, err := wH.wS.FindByCustomerID(c.Context(), customerID)
	if err != nil {
		// 500 server error finding by customer id
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 200 ok
	return c.Status(fiber.StatusOK).JSON(workOrders)
}
