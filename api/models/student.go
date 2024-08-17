package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedBy uuid.UUID `json:"updated_by"`
	UpdatedOn time.Time `json:"updated_on"`
}
