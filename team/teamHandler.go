package team

import "net/http"

type TeamHandler interface {
	AddToProject(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveFromProject(w http.ResponseWriter, r *http.Request) (int, error)
}

type TeamHandlerImpl struct {
	storage TeamStorage
	parser  TeamParser
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

func NewTeamHandler(storage TeamStorage) TeamHandler {
	return &TeamHandlerImpl{
		storage: storage,
		parser:  newTeamParser(),
	}
}