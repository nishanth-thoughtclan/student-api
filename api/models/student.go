package models

import "time"

type Student struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedBy string    `json:"created_by"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedBy string    `json:"updated_by"`
	UpdatedOn time.Time `json:"updated_on"`
}
