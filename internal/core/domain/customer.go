// internal/core/domain/customer.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
	Address   string    `gorm:"not null"`
	//puntero para poder capturar el nil
	StartDate  *time.Time
	EndDate    *time.Time
	IsActive   bool      `gorm:"not null; default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	WorkOrders []WorkOrder
}
