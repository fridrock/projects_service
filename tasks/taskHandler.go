package tasks

import (
	"encoding/json"
	"net/http"
)

type TaskHandler interface {
	AddToBacklog(w http.ResponseWriter, r *http.Request) (int, error)
	SetExecutor(w http.ResponseWriter, r *http.Request) (int, error)
	SetColumn(w http.ResponseWriter, r *http.Request) (int, error)
	DeleteTask(w http.ResponseWriter, r *http.Request) (int, error)
	GetProjectTasks(w http.ResponseWriter, r *http.Request) (int, error)
}

type TaskHandlerImpl struct {
	storage TaskStorage
	parser  TaskParser
}

func (th *TaskHandlerImpl) AddToBacklog(w http.ResponseWriter, r *http.Request) (int, error) {
	task, err := th.parser.GetTask(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	task, err = th.storage.AddToBacklog(task)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	responseText, err := json.MarshalIndent(task, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}

func (th *TaskHandlerImpl) SetExecutor(w http.ResponseWriter, r *http.Request) (int, error) {
	columnDto, err := th.parser.GetExecutorDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = th.storage.SetExecutor(columnDto.ExecutorId, columnDto.TaskId)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func (th *TaskHandlerImpl) SetColumn(w http.ResponseWriter, r *http.Request) (int, error) {
	columnDto, err := th.parser.GetColumnDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = th.storage.SetColumn(columnDto.ColumnId, columnDto.TaskId)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func (th *TaskHandlerImpl) DeleteTask(w http.ResponseWriter, r *http.Request) (int, error) {
	taskId, err := th.parser.GetDeleteTask(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = th.storage.DeleteTask(taskId)
	if err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
func (th *TaskHandlerImpl) GetProjectTasks(w http.ResponseWriter, r *http.Request) (int, error) {
	dto, err := th.parser.GetProjectTasksDto(r)
	if err != nil {
		return http.StatusBadRequest, err
	}
	tasks, err := th.storage.GetProjectTasks(dto.ProjectId)
	if err != nil {
		return http.StatusNotFound, err
	}
	responseText, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseText)
	return http.StatusOK, nil
}

func NewTaskHandler(storage TaskStorage) TaskHandler {
	return &TaskHandlerImpl{
		storage: storage,
		parser:  newTaskParser(),
	}
}
