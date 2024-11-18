package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/fridrock/projects_service/db/core"
	"github.com/fridrock/projects_service/project"
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
	mainRouter.Handle("/projects/", utils.HandleErrorMiddleware(a.authManager.HandleWithAuth(a.projectHandler.DeleteProject))).Methods("DELETE")
	return mainRouter
}
