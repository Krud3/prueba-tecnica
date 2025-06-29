// internal/core/ports/ports.go

package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/krud3/prueba-tecnica/internal/core/domain"
)

type WorkOrderFilters struct {
	Since  *time.Time
	Until  *time.Time
	Status *domain.Status
}

type CustomerRepository interface {
	Create(ctx context.Context, customer domain.Customer) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	GetActive(ctx context.Context) ([]domain.Customer, error)
	Update(ctx context.Context, customer domain.Customer) error
}

type WorkOrderRepository interface {
	Create(ctx context.Context, workOrder domain.WorkOrder) error
	CompleteOrder(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkOrder, error)
	FindByFilter(ctx context.Context, filters WorkOrderFilters) ([]domain.WorkOrder, error)
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.WorkOrder, error)
	Update(ctx context.Context, workOrder domain.WorkOrder) error
}
