// internal/core/domain/workorder.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status string

type Type string

const (
	StatusNew       Status = "new"
	StatusDone      Status = "done"
	StatusCancelled Status = "cancelled"
)

const (
	TypeActivate Type = "activar cliente"
	TypeCancell  Type = "cancelar cliente"
)

type WorkOrder struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	CustomerID       uuid.UUID `gorm:"type:uuid;not null"`
	Customer         Customer  `gorm:"foreignKey:CustomerID"` //facilita la condicion 9
	Description      string    `gorm:"not null"`
	PlannedDateBegin time.Time `gorm:"not null"`
	PlannedDateEnd   time.Time `gorm:"not null"`
	Status           Status    `gorm:"type:enum('new','done','cancelled');default:'new';not null"`
	Type             Type      `gorm:"not null"` //debido a la logica de negocio, definimos a type como dos valores, pero realmente a gorm le mandamos un string, si se maneja con un enum o un default, podria causar errores en el futuro
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}
