// internal/core/services/services.go

package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/krud3/prueba-tecnica/internal/core/domain"
	"github.com/krud3/prueba-tecnica/internal/core/ports"
)

var (
	// handle error for active customers that recieve active orders
	ErrAA = errors.New("el cliente ya tiene su estado activado")

	// handle error for no active customers that recieve cancel order
	ErrCC = errors.New("el cliente ya tiene su estado cancelado")

	// handle error for planned date interval
	ErrDateIntertal = errors.New("la diferencia entre las fechas de planeación no debe ser mayor a dos horas")

	// handle error for workOrder done trying to be completed
	ErrWODone = errors.New("la orden ya estaba completada y está intentando completarla")

	// handle error for workOrder canceled trying to be completed
	ErrWOCancelled = errors.New("la orden está cancelada y está intentando completarla")
)

type CustomerService struct {
	cRepo ports.CustomerRepository
}

func NewCustomerService(customerRepo ports.CustomerRepository) *CustomerService {
	return &CustomerService{cRepo: customerRepo}
}

func (cS *CustomerService) Create(ctx context.Context, customer domain.Customer) error {
	return cS.cRepo.Create(ctx, customer)
}

func (cS *CustomerService) FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	return cS.cRepo.FindByID(ctx, id)
}

func (cS *CustomerService) GetActive(ctx context.Context) ([]domain.Customer, error) {
	return cS.cRepo.GetActive(ctx)
}

func (cS *CustomerService) GetAll(ctx context.Context) ([]domain.Customer, error) {
	return cS.cRepo.GetAll(ctx)
}

func (cS *CustomerService) Update(ctx context.Context, customer domain.Customer) error {
	return cS.cRepo.Update(ctx, customer)
}

type WorkOrderService struct {
	wRepo      ports.WorkOrderRepository
	cRepo      ports.CustomerRepository
	redis      *redis.Client
	streamName string
}

func NewWorkOrderService(workOrderRepo ports.WorkOrderRepository, customerRepo ports.CustomerRepository, redisClient *redis.Client, stream string) *WorkOrderService {
	return &WorkOrderService{
		wRepo:      workOrderRepo,
		cRepo:      customerRepo,
		redis:      redisClient,
		streamName: stream,
	}
}

func (wS *WorkOrderService) Create(ctx context.Context, workOrder domain.WorkOrder) error {

	// compares end and begin not > 2 #business logic 2
	if workOrder.PlannedDateEnd.Sub(workOrder.PlannedDateBegin).Hours() > 2 {
		return ErrDateIntertal
	}

	// get customer
	customer, err := wS.cRepo.FindByID(ctx, workOrder.CustomerID)
	if err != nil {
		return err
	}

	// switch customer.IsActive from domain to handle errors
	switch customer.IsActive {
	case true:
		// trying to activate a customer already active
		if workOrder.Type == domain.TypeActivate {
			return ErrAA
		}
	case false:
		// trying to cancel a customer notActive
		if workOrder.Type == domain.TypeCancell {
			return ErrCC
		}
	}

	// create workOrder
	return wS.wRepo.Create(ctx, workOrder)
}

// handles CompleteOrder for business conditions
func (wS *WorkOrderService) CompleteOrder(ctx context.Context, id uuid.UUID) error {
	// check if workOrder exist by ID
	workOrder, err := wS.wRepo.FindByID(ctx, id)
	// handle error
	if err != nil {
		return err
	}

	switch workOrder.Status {
	// handles workOrder.Status not been hable to change to Done while Done at current status
	case domain.StatusDone:
		return ErrWODone
	// handles workOrder.Status not been hable to change to Done while Cancelled at current status
	case domain.StatusCancelled:
		return ErrWOCancelled
	}

	// check if customer exist by ID given by workOrder struct
	customer, err := wS.cRepo.FindByID(ctx, workOrder.CustomerID)
	// handle error
	if err != nil {
		return err
	}
	// time for set StartDate or EndDate
	timeNow := time.Now()

	// set isActive to costumer
	switch workOrder.Type {
	case domain.TypeActivate:
		customer.IsActive = true
		customer.StartDate = &timeNow
		customer.EndDate = nil

	case domain.TypeCancell:
		customer.IsActive = false
		customer.EndDate = &timeNow
	}

	// make the change doing Update passing customer pointer
	errUC := wS.cRepo.Update(ctx, *customer)
	if errUC != nil {
		return errUC
	}

	// set Status to workOrder
	workOrder.Status = domain.StatusDone
	// make the change to workOrder passing workOrder pointer
	errUWO := wS.wRepo.Update(ctx, *workOrder)
	if errUWO != nil {
		return errUWO
	}

	// map workOrder into json to send it to redis
	workOrderJSON, err := json.Marshal(workOrder)
	if err != nil {
		return err
	}
	// send the event to redis stream
	err = wS.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: wS.streamName,
		Values: map[string]interface{}{
			"work_order_completed": string(workOrderJSON),
		},
	}).Err()

	if err != nil {
		return err
	}

	return nil
}

func (wS *WorkOrderService) FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkOrder, error) {
	return wS.wRepo.FindByID(ctx, id)
}

func (wS *WorkOrderService) FindByFilter(ctx context.Context, filters ports.WorkOrderFilters) ([]domain.WorkOrder, error) {
	return wS.wRepo.FindByFilter(ctx, filters)
}

func (wS *WorkOrderService) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.WorkOrder, error) {
	return wS.wRepo.FindByCustomerID(ctx, customerID)
}

func (wS *WorkOrderService) Update(ctx context.Context, workOrder domain.WorkOrder) error {
	return wS.wRepo.Update(ctx, workOrder)
}
