package columns

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fridrock/projects_service/api"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TeammateDto struct {
	UserId    uuid.UUID `json:"userId" validate:"required"`
	ProjectId uuid.UUID `json:"projectId" validate:"required"`
}
type ColumnByProjectDto struct {
	ProjectId uuid.UUID `json:"projectId" validate:"required"`
}

type ColumnParser interface {
	GetColumn(r *http.Request) (api.Column, error)
	GetDeleteColumnId(r *http.Request) (uuid.UUID, error)
	GetColumnByProject(r *http.Request) (uuid.UUID, error)
}

type ColumnParserImpl struct {
	validate *validator.Validate
}

func (cp *ColumnParserImpl) GetColumn(r *http.Request) (api.Column, error) {
	var dto api.Column
	err := json.NewDecoder(r.Body).Decode(&dto)
	return dto, errors.Join(err, cp.validate.Struct(dto))
}

func (cp *ColumnParserImpl) GetDeleteColumnId(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	columnId, err := uuid.Parse(vars["id"])
	return columnId, err
}

func (cp *ColumnParserImpl) GetColumnByProject(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	projectId, err := uuid.Parse(vars["projectId"])
	return projectId, err
}

func newTeamParser() ColumnParser {
	return &ColumnParserImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}
