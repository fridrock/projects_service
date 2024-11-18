package api

import "github.com/google/uuid"

type Project struct {
	Id      uuid.UUID `json:"id" db:"id"`
	Name    string    `json:"name" db:"name" validate:"required"`
	OwnerId uuid.UUID `json:"owner_id" db:"owner_id" validate:"required"`
}
