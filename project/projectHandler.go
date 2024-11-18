package project

import (
	"encoding/json"
	"net/http"

	"github.com/fridrock/projects_service/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProjectHandler interface {
	CreateProject(w http.ResponseWriter, r *http.Request) (int, error)
	GetProjects(w http.ResponseWriter, r *http.Request) (int, error)
	DeleteProject(w http.ResponseWriter, r *http.Request) (int, error)
}

type ProjectHandlerImpl struct {
	storage ProjectStorage
	parser  ProjectParser
}

func (ph *ProjectHandlerImpl) CreateProject(w http.ResponseWriter, r *http.Request) (int, error) {
	projectDto, err := ph.parser.GetProjectDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	projectSaved, err := ph.storage.CreateProject(projectDto)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	responseText, err := json.MarshalIndent(projectSaved, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}

func (ph *ProjectHandlerImpl) GetProjects(w http.ResponseWriter, r *http.Request) (int, error) {
	userId := utils.UserFromContext(r.Context())
	projects, err := ph.storage.GetProjects(userId)
	if err != nil {
		return http.StatusNotFound, err
	}
	responseText, err := json.MarshalIndent(projects, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}

func (ph *ProjectHandlerImpl) DeleteProject(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	var projectId uuid.UUID
	projectId, err := uuid.Parse(vars["id"])
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = ph.storage.DeleteProject(projectId)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func NewProjectHandler(storage ProjectStorage) ProjectHandler {
	return &ProjectHandlerImpl{
		storage: storage,
		parser:  newProjectParser(),
	}
}
