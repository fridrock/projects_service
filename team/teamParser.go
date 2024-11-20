package team

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TeammateDto struct {
	UserId    uuid.UUID `json:"userId" validate:"required"`
	ProjectId uuid.UUID `json:"projectId" validate:"required"`
}

type TeamParser interface {
	ParseTeammateDto(r *http.Request) (TeammateDto, error)
	ParseProjectId(r *http.Request) (uuid.UUID, error)
}

type TeamParserImpl struct {
	validate *validator.Validate
}

func (tp *TeamParserImpl) ParseTeammateDto(r *http.Request) (TeammateDto, error) {
	var dto TeammateDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	return dto, errors.Join(err, tp.validate.Struct(dto))
}

func (tp *TeamParserImpl) ParseProjectId(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["projectId"])
}

func newTeamParser() TeamParser {
	return &TeamParserImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}
