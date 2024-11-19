package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/fridrock/projects_service/db/core"
	"github.com/fridrock/projects_service/project"
	"github.com/fridrock/projects_service/tasks"
	"github.com/fridrock/projects_service/team"
	"github.com/fridrock/projects_service/utils"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	startApp()
}

type App struct {
	server         *http.Server
	db             *sqlx.DB
	projectStorage project.ProjectStorage
	projectHandler project.ProjectHandler
	teamStorage    team.TeamStorage
	teamHandler    team.TeamHandler
	taskStorage    tasks.TaskStorage
	taskHandler    tasks.TaskHandler
	authManager    utils.AuthManager
}

func startApp() {
	a := App{}
	a.setup()
}
func (a App) setup() {
	a.db = core.CreateConnection()
	defer a.db.Close()
	a.projectStorage = project.NewProjectStorage(a.db)
	a.projectHandler = project.NewProjectHandler(a.projectStorage)
	a.teamStorage = team.NewTeamStorage(a.db)
	a.teamHandler = team.NewTeamHandler(a.teamStorage)
	a.taskStorage = tasks.NewTaskStorage(a.db)
	a.taskHandler = tasks.NewTaskHandler(a.taskStorage)
	a.authManager = utils.NewAuthManager()
	a.setupServer()
}

func (a App) setupServer() {
	a.server = &http.Server{
		Addr:         ":8000",
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Handler:      a.getRouter(),
	}
	slog.Info("Starting server on port 8000")
	a.server.ListenAndServe()
}
func (a App) getRouter() http.Handler {
	mainRouter := mux.NewRouter()
	mainRouter.Handle("/projects/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.projectHandler.CreateProject))).Methods("POST")
	mainRouter.Handle("/projects/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.projectHandler.GetProjects))).Methods("GET")
	mainRouter.Handle("/projects/{id}", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.projectHandler.DeleteProject))).Methods("DELETE")
	mainRouter.Handle("/team/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.teamHandler.AddToProject))).Methods("POST")
	mainRouter.Handle("/team/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.teamHandler.RemoveFromProject))).Methods("DELETE")
	mainRouter.Handle("/task/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.GetProjectTasks))).Methods("GET")
	mainRouter.Handle("/task/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.AddToBacklog))).Methods("POST")
	mainRouter.Handle("/task/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.SetExecutor))).Methods("PATCH")
	mainRouter.Handle("/task/{id}", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.DeleteTask))).Methods("DELETE")
	return mainRouter
}
