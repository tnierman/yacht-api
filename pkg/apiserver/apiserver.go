package apiserver

import (
	"net/http"
	"time"

	"github.com/tnierman/yacht-api/pkg/logging"
	"github.com/tnierman/yacht-api/pkg/services/jira"
	"github.com/tnierman/yacht-api/pkg/services/pagerduty"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
	server *http.Server

	// Services
	Pagerduty pagerduty.Service
	Jira jira.Service
}

func NewServer() Server {
	router := mux.NewRouter()
	pagerdutyService := pagerduty.NewService()
	jiraService := jira.NewService()

	pagerdutyRouter := router.PathPrefix("/pagerduty").Subrouter()
	pagerdutyService.AddRoutes(pagerdutyRouter)

	jiraRouter := router.PathPrefix("/jira").Subrouter()
	jiraService.AddRoutes(jiraRouter)


	s := &http.Server{
		Handler: router,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	server := Server{
		Router: router,
		server: s,
		Pagerduty: pagerdutyService,
		Jira: jiraService,
	}
	return server
}

func (s *Server) Serve(logger *logging.Logger) error {
	logger.Info("todo")

	return s.server.ListenAndServe()
}
