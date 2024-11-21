package tasks

import (
	"log/slog"

	"github.com/fridrock/projects_service/api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskStorage interface {
	AddToBacklog(api.Task) (api.Task, error)
	DeleteTask(uuid.UUID) error
	SetExecutor(executorId uuid.UUID, taskId uuid.UUID) error
	SetColumn(columnId uuid.UUID, taskId uuid.UUID) error
	GetProjectTasks(uuid.UUID) ([]api.Task, error)
}
type TaskStorageImpl struct {
	db *sqlx.DB
}

func (ts *TaskStorageImpl) AddToBacklog(task api.Task) (api.Task, error) {
	getLastNumQ := `SELECT num FROM tasks WHERE project_id = $1 ORDER BY num DESC LIMIT 1`
	var lastNum int
	err := ts.db.Get(&lastNum, getLastNumQ, task.ProjectId)
	if err != nil {
		slog.Debug(err.Error())
	}
	lastNum += 1
	q := `INSERT INTO tasks(id, num, name,  description, project_id) VALUES($1, $2, $3, $4, $5) RETURNING id`
	err = ts.db.QueryRow(q, uuid.New(), lastNum, task.Name, task.Description, task.ProjectId).Scan(&task.Id)
	task.Num = lastNum
	return task, err
}
func (ts *TaskStorageImpl) DeleteTask(id uuid.UUID) error {
	q := `DELETE FROM tasks WHERE id = $1`
	_, err := ts.db.Exec(q, id)
	return err
}

func (ts *TaskStorageImpl) SetExecutor(executorId uuid.UUID, taskId uuid.UUID) error {
	q := `UPDATE tasks SET executor_id = $1 WHERE id = $2`
	_, err := ts.db.Exec(q, executorId, taskId)
	return err
}

func (ts *TaskStorageImpl) SetColumn(columnId uuid.UUID, taskId uuid.UUID) error {
	q := `UPDATE tasks SET column_id = $1 WHERE id = $2`
	_, err := ts.db.Exec(q, columnId, taskId)
	return err
}

func (ts *TaskStorageImpl) GetProjectTasks(projectId uuid.UUID) ([]api.Task, error) {
	var tasks []api.Task
	q := `SELECT * FROM tasks WHERE project_id = $1`
	err := ts.db.Select(&tasks, q, projectId)
	return tasks, err
}

func NewTaskStorage(db *sqlx.DB) TaskStorage {
	return &TaskStorageImpl{
		db: db,
	}
}
