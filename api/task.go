package api

import "github.com/google/uuid"

type Task struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Num         int       `json:"num" db:"num"`
	Description string    `json:"description" db:"description" validate:"required"`
	ProjectId   uuid.UUID `json:"projectId" db:"project_id" validate:"required"`
	ExecutorId  uuid.UUID `json:"executorId" db:"executor_id"`
	ColumnId    uuid.UUID `json:"columnId" db:"column_id"`
}
