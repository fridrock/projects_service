package team

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeamStorage interface {
	AddToProject(userId uuid.UUID, projectId uuid.UUID) error
	RemoveFromProject(userId uuid.UUID, projectId uuid.UUID) error
	GetTeammatesIds(projectId uuid.UUID) ([]uuid.UUID, error)
}
type TeamStorageImpl struct {
	db *sqlx.DB
}

func (ts *TeamStorageImpl) AddToProject(userId uuid.UUID, projectId uuid.UUID) error {
	q := `INSERT INTO users_in_projects(user_id, project_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := ts.db.Exec(q, userId, projectId)
	return err
}
func (ts *TeamStorageImpl) RemoveFromProject(userId uuid.UUID, projectId uuid.UUID) error {
	q := `DELETE FROM users_in_projects WHERE user_id = $1 AND project_id = $2`
	_, err := ts.db.Exec(q, userId, projectId)
	return err
}
func (ts *TeamStorageImpl) GetTeammatesIds(projectId uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	q := `SELECT user_id FROM users_in_projects WHERE project_id = $1`
	err := ts.db.Select(&ids, q, projectId)
	return ids, err
}

func NewTeamStorage(db *sqlx.DB) TeamStorage {
	return &TeamStorageImpl{
		db: db,
	}
}
