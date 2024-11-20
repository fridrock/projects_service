package team

import (
	"encoding/json"
	"net/http"

	"github.com/fridrock/projects_service/api"
)

type TeamHandler interface {
	AddToProject(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveFromProject(w http.ResponseWriter, r *http.Request) (int, error)
	GetProfiles(w http.ResponseWriter, r *http.Request) (int, error)
}

type TeamHandlerImpl struct {
	usersClient UsersClient
	storage     TeamStorage
	parser      TeamParser
}

func (ts *TeamHandlerImpl) AddToProject(w http.ResponseWriter, r *http.Request) (int, error) {
	dto, err := ts.parser.ParseTeammateDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = ts.storage.AddToProject(dto.UserId, dto.ProjectId)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (ts *TeamHandlerImpl) RemoveFromProject(w http.ResponseWriter, r *http.Request) (int, error) {
	dto, err := ts.parser.ParseTeammateDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = ts.storage.RemoveFromProject(dto.UserId, dto.ProjectId)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
func (ts *TeamHandlerImpl) GetProfiles(w http.ResponseWriter, r *http.Request) (int, error) {
	projectId, err := ts.parser.ParseProjectId(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	ids, err := ts.storage.GetTeammatesIds(projectId)
	if err != nil {
		return http.StatusNotFound, err
	}
	profiles, err := ts.usersClient.GetProfilesByIds(api.ProfilesByUserIdsDto{
		Ids: ids,
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	responseText, err := json.MarshalIndent(profiles, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}
func NewTeamHandler(storage TeamStorage) TeamHandler {
	return &TeamHandlerImpl{
		storage:     storage,
		usersClient: NewUsersClient(),
		parser:      newTeamParser(),
	}
}
