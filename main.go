package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/fridrock/projects_service/columns"
	"github.com/fridrock/projects_service/db/core"
	"github.com/fridrock/projects_service/project"
	"github.com/fridrock/projects_service/tasks"
	"github.com/fridrock/projects_service/team"
	"github.com/fridrock/projects_service/utils"
	"github.com/gorilla/handlers"
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
	columnStorage  columns.ColumnStorage
	columnHandler  columns.ColumnHandler
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
	a.columnStorage = columns.NewColumnStorage(a.db)
	a.columnHandler = columns.NewColumnHandler(a.columnStorage)
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
	mainRouter.Handle("/projects/{id}", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.projectHandler.DeleteProject))).Methods("DELETE", "OPTIONS")
	mainRouter.Handle("/team/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.teamHandler.AddToProject))).Methods("POST")
	mainRouter.Handle("/team/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.teamHandler.RemoveFromProject))).Methods("DELETE")
	mainRouter.Handle("/team/profiles/{projectId}", utils.HandleErrorMiddleware(a.teamHandler.GetProfiles)).Methods("GET")
	mainRouter.Handle("/task/byproject/{projectId}", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.GetProjectTasks))).Methods("GET")
	mainRouter.Handle("/task/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.AddToBacklog))).Methods("POST")
	mainRouter.Handle("/task/executor", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.SetExecutor))).Methods("PATCH")
	mainRouter.Handle("/task/column", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.SetColumn))).Methods("PATCH")
	mainRouter.Handle("/task/{id}", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.taskHandler.DeleteTask))).Methods("DELETE")
	mainRouter.Handle("/column/byproject/{projectId}", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.columnHandler.GetColumnByProject))).Methods("GET")
	mainRouter.Handle("/column/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.columnHandler.AddToProject))).Methods("POST")
	mainRouter.Handle("/column/{id}", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.columnHandler.RemoveFromProject))).Methods("DELETE")
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	handler := corsObj(mainRouter)
	return handler
}
