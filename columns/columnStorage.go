package columns

import (
	"github.com/fridrock/projects_service/api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ColumnStorage interface {
	AddToProject(api.Column) (api.Column, error)
	RemoveFromProject(uuid.UUID) error
	GetColumnsByProject(uuid.UUID) ([]api.ColumnExtended, error)
}
type ColumnStorageImpl struct {
	db *sqlx.DB
}

func (cs *ColumnStorageImpl) AddToProject(dto api.Column) (api.Column, error) {
	q := `INSERT INTO project_columns(id, name, project_id) VALUES ($1, $2, $3) RETURNING id`
	err := cs.db.QueryRow(q, uuid.New(), dto.Name, dto.ProjectId).Scan(&dto.Id)
	return dto, err
}
func (cs *ColumnStorageImpl) RemoveFromProject(columnId uuid.UUID) error {
	q := `DELETE FROM project_columns WHERE id = $1`
	_, err := cs.db.Exec(q, columnId)
	return err
}

func (cs *ColumnStorageImpl) GetColumnsByProject(projectId uuid.UUID) ([]api.ColumnExtended, error) {
	var columns []api.ColumnExtended
	q := `SELECT * FROM project_columns WHERE project_id = $1`
	q2 := `SELECT * FROM tasks WHERE column_id = $1`
	err := cs.db.Select(&columns, q, projectId)
	if err != nil {
		return columns, err
	}
	for i := 0; i < len(columns); i++ {
		var tasks []api.Task
		err = cs.db.Select(&tasks, q2, columns[i].Id)
		if err != nil {
			return columns, err
		}
		columns[i].Tasks = tasks
	}
	return columns, nil
}

func NewColumnStorage(db *sqlx.DB) ColumnStorage {
	return &ColumnStorageImpl{
		db: db,
	}
}
