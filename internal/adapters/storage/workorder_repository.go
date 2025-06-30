// internal/adapters/storage/workorder_repository.go

package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/krud3/prueba-tecnica/internal/core/domain"
	"github.com/krud3/prueba-tecnica/internal/core/ports"
	"gorm.io/gorm"
)

var (
	ErrNoWID = errors.New("no se encontró ID asociada al work order")
)

type gormWorkOrderRepository struct {
	db *gorm.DB
}

func NewGormWorkOrderRepository(db *gorm.DB) ports.WorkOrderRepository {
	return &gormWorkOrderRepository{db: db}
}

func (r *gormWorkOrderRepository) Create(ctx context.Context, workOrder domain.WorkOrder) error {
	if workOrder.ID == uuid.Nil {
		workOrder.ID = uuid.New()
	}
	// SQL Insert with create
	result := r.db.WithContext(ctx).Create(&workOrder)

	// if there is any error return it
	return result.Error
}

func (r *gormWorkOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.WorkOrder, error) {
	var workOrder domain.WorkOrder

	// preload of customer since condition 9 especifys it, due workOrder pointer SELECT assigns workOrder finded to var workOrder
	result := r.db.WithContext(ctx).Preload("Customer").First(&workOrder, "id = ?", id)
	if result.Error != nil {
		// if not found nil nil
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// if error
		return nil, result.Error
	}
	// founded
	return &workOrder, nil
}

func (r *gormWorkOrderRepository) FindByFilter(ctx context.Context, filters ports.WorkOrderFilters) ([]domain.WorkOrder, error) {
	// stores workOrders finded if any
	var workOrders []domain.WorkOrder
	// specify wich talbe gorm is working on
	query := r.db.WithContext(ctx).Model(&domain.WorkOrder{})

	// joins where according to filters
	if filters.Since != nil {
		query = query.Where("planned_date_begin >= ?", *filters.Since)
	}
	if filters.Until != nil {
		query = query.Where("planned_date_end <= ?", *filters.Until)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	// Preload customer and storage results in workOrders
	err := query.Preload("Customer").Find(&workOrders).Error
	if err != nil {
		return nil, err
	}

	return workOrders, nil
}

func (r *gormWorkOrderRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.WorkOrder, error) {
	// stores workOrders if any
	var workOrders []domain.WorkOrder

	// where filter Preloads Customer storage in workOrders
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Preload("Customer").Find(&workOrders).Error

	return workOrders, err
}

func (r *gormWorkOrderRepository) Update(ctx context.Context, workOrder domain.WorkOrder) error {
	// if no id given error
	if workOrder.ID == uuid.Nil {
		return ErrNoWID
	} else {
		// save update value in db
		result := r.db.WithContext(ctx).Save(&workOrder)
		return result.Error
	}

}
