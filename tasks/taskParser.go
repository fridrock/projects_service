package tasks

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fridrock/projects_service/api"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SetExecutorDto struct {
	ExecutorId uuid.UUID `json:"executorId" validate:"required"`
	TaskId     uuid.UUID `json:"taskId" validate:"required"`
}
type ProjectTasksDto struct {
	ProjectId uuid.UUID `json:"projectId" validate:"required"`
}

type TaskParser interface {
	GetTask(r *http.Request) (api.Task, error)
	GetExecutorDto(r *http.Request) (SetExecutorDto, error)
	GetDeleteTask(r *http.Request) (taskId uuid.UUID, err error)
	GetProjectTasksDto(r *http.Request) (ProjectTasksDto, error)
}

type TaskParserImpl struct {
	validate *validator.Validate
}

func (tp *TaskParserImpl) GetTask(r *http.Request) (api.Task, error) {
	var dto api.Task
	err := json.NewDecoder(r.Body).Decode(&dto)
	err = errors.Join(err, tp.validate.Struct(dto))
	return dto, err
}

func (tp *TaskParserImpl) GetExecutorDto(r *http.Request) (SetExecutorDto, error) {
	var dto SetExecutorDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	err = errors.Join(err, tp.validate.Struct(dto))
	return dto, err
}

func (tp *TaskParserImpl) GetDeleteTask(r *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(r)
	taskId, err := uuid.Parse(vars["id"])
	return taskId, err
}
func (tp *TaskParserImpl) GetProjectTasksDto(r *http.Request) (ProjectTasksDto, error) {
	var dto ProjectTasksDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	err = errors.Join(err, tp.validate.Struct(dto))
	return dto, err
}

func newTaskParser() TaskParser {
	return &TaskParserImpl{
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}
