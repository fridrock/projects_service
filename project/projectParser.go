package project

import (
	"encoding/json"
	"net/http"

	"github.com/fridrock/projects_service/api"
	"github.com/fridrock/projects_service/utils"
	"github.com/go-playground/validator/v10"
)

type ProjectParser interface {
	GetProjectDto(*http.Request) (api.Project, error)
}

type ProjectParserImpl struct {
	validate *validator.Validate
}

func newProjectParser() ProjectParser {
	return ProjectParserImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (pp ProjectParserImpl) GetProjectDto(r *http.Request) (api.Project, error) {
	var projectDto api.Project
	err := json.NewDecoder(r.Body).Decode(&projectDto)
	if err != nil {
		return projectDto, err
	}
	projectDto.OwnerId = utils.UserFromContext(r.Context())
	err = pp.validate.Struct(projectDto)
	return projectDto, err
}
