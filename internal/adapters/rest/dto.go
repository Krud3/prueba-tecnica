// internal/adapters/rest/deteo.go

package rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/krud3/prueba-tecnica/internal/core/domain"
)

type CreateCustomerRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
}

type CreateWorkOrderRequest struct {
	CustomerID       uuid.UUID   `json:"customerID"`
	Description      string      `json:"description"`
	PlannedDateBegin time.Time   `json:"plannedDateBegin"`
	PlannedDateEnd   time.Time   `json:"plannedDateEnd"`
	Type             domain.Type `json:"type"`
}
