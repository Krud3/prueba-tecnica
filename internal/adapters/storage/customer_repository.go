// internal/adapters/storage/customer_repository.go

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
	ErrNoCID = errors.New("no se encontró ID asociada al customer")
)

type gormCustomerRepository struct {
	db *gorm.DB
}

func NewGormCustomerRepository(db *gorm.DB) ports.CustomerRepository {
	return &gormCustomerRepository{db: db}
}

func (r *gormCustomerRepository) Create(ctx context.Context, customer domain.Customer) error {
	//uuid if not exist
	if customer.ID == uuid.Nil {
		customer.ID = uuid.New()
	}
	// create customer
	result := r.db.WithContext(ctx).Create(&customer)

	// if any error return else nil
	return result.Error
}

func (r *gormCustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	// to storage customer
	var customer domain.Customer

	// assings customer by pointer
	result := r.db.WithContext(ctx).First(&customer, "id = ?", id)
	if result.Error != nil {
		// if error
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// if not found
		return nil, result.Error
	}
	// founded
	return &customer, nil
}

func (r *gormCustomerRepository) GetActive(ctx context.Context) ([]domain.Customer, error) {
	// to storages customers
	var customers []domain.Customer

	// search is active and stores it catch error if error
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&customers).Error

	// results
	return customers, err
}

func (r *gormCustomerRepository) GetAll(ctx context.Context) ([]domain.Customer, error) {
	var customers []domain.Customer

	err := r.db.WithContext(ctx).Find(&customers).Error

	return customers, err
}

func (r *gormCustomerRepository) Update(ctx context.Context, customer domain.Customer) error {
	if customer.ID == uuid.Nil {
		return ErrNoCID
	} else {
		return r.db.WithContext(ctx).Save(customer).Error
	}
}
