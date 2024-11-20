package api

import "github.com/google/uuid"

type Column struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required"`
	ProjectId uuid.UUID `json:"projectId" db:"project_id" validate:"required"`
}

type ColumnExtended struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ProjectId uuid.UUID `json:"projectId" db:"project_id"`
	Tasks     []Task    `json:"tasks"`
}
