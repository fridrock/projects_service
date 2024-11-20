package columns

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ColumnHandler interface {
	AddToProject(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveFromProject(w http.ResponseWriter, r *http.Request) (int, error)
	GetColumnByProject(w http.ResponseWriter, r *http.Request) (int, error)
}

type ColumnHandlerImpl struct {
	storage ColumnStorage
	parser  ColumnParser
}

func (ch *ColumnHandlerImpl) AddToProject(w http.ResponseWriter, r *http.Request) (int, error) {
	dto, err := ch.parser.GetColumn(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	dto, err = ch.storage.AddToProject(dto)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	responseText, err := json.MarshalIndent(dto, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}

func (ch *ColumnHandlerImpl) RemoveFromProject(w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := ch.parser.GetDeleteColumnId(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = ch.storage.RemoveFromProject(id)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func (ch *ColumnHandlerImpl) GetColumnByProject(w http.ResponseWriter, r *http.Request) (int, error) {
	projectId, err := ch.parser.GetColumnByProject(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	columns, err := ch.storage.GetColumnsByProject(projectId)
	if err != nil {
		slog.Debug(err.Error())
		return http.StatusNotFound, err
	}
	responseText, err := json.MarshalIndent(columns, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}

func NewColumnHandler(storage ColumnStorage) ColumnHandler {
	return &ColumnHandlerImpl{
		storage: storage,
		parser:  newTeamParser(),
	}
}
