package project

import (
	"github.com/fridrock/projects_service/api"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProjectStorage interface {
	CreateProject(api.Project) (api.Project, error)
	GetProjects(userId uuid.UUID) ([]api.Project, error)
	DeleteProject(projectId uuid.UUID) error
	GetProject(userId uuid.UUID, projectId uuid.UUID) (api.Project, error)
}

type ProjectStorageImpl struct {
	db *sqlx.DB
}

func (ps *ProjectStorageImpl) CreateProject(projectDto api.Project) (api.Project, error) {
	q := `INSERT INTO projects(id, name, owner_id) VALUES ($1,$2, $3) RETURNING id`
	err := ps.db.QueryRow(q, uuid.New(), projectDto.Name, projectDto.OwnerId).Scan(&projectDto.Id)
	if err != nil {
		return projectDto, err
	}
	q = `INSERT INTO users_in_projects(user_id, project_id) VALUES ($1, $2)`
	_, err = ps.db.Exec(q, projectDto.OwnerId, projectDto.Id)
	return projectDto, err
}

func (ps *ProjectStorageImpl) GetProjects(userId uuid.UUID) ([]api.Project, error) {
	var projectList []api.Project
	q := `SELECT projects.id, projects.name, projects.owner_id FROM users_in_projects LEFT JOIN projects ON projects.id = project_id WHERE user_id = $1`
	err := ps.db.Select(&projectList, q, userId)
	return projectList, err
}

func (ps *ProjectStorageImpl) DeleteProject(projectId uuid.UUID) error {
	q := `DELETE FROM projects WHERE id = $1`
	_, err := ps.db.Exec(q, projectId)
	return err
}

func NewProjectStorage(db *sqlx.DB) ProjectStorage {
	return &ProjectStorageImpl{
		db: db,
	}
}
