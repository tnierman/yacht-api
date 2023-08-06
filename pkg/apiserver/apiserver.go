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
	logger *logging.Logger

	// Services
	Pagerduty *pagerduty.Service
	Jira jira.Service
}

func NewServer(parent *logging.Logger, opts ServerOptions) Server {
	// Initialize server internals
	router := mux.NewRouter()
	logger := parent.NewChildLogger("server")

	// Initialize services
	pagerdutyService := pagerduty.NewService(logger, opts.Pagerduty)
	pagerdutyRouter := router.PathPrefix("/pagerduty").Subrouter()
	pagerdutyService.AddRoutes(pagerdutyRouter)

	jiraService := jira.NewService()
	jiraRouter := router.PathPrefix("/jira").Subrouter()
	jiraService.AddRoutes(jiraRouter)

	// Build server from components
	s := &http.Server{
		Handler: router,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	server := Server{
		// Internals
		Router: router,
		server: s,
		logger: logger,

		// Services
		Pagerduty: pagerdutyService,
		Jira: jiraService,
	}
	return server
}

// Serve listens for and serves requests
func (s *Server) Serve(logger *logging.Logger) error {
	logger.Info("todo")

	// Start services
	s.Pagerduty.Start(logger)

	defer func() {
		cleanupErr := s.Cleanup()
		if cleanupErr != nil {
			logger.Sugar().Errorf("failed to properly cleanup server: %v", cleanupErr)
		}
	}()
	err := s.server.ListenAndServe()
	logger.Sugar().Errorf("stopping server, error while serving: %v", err)
	return err
}

// Cleanup stops this Server's services and closes the server
func (s *Server) Cleanup() error {
	// Stop services
	s.Pagerduty.Stop()

	// Stop server
	return s.server.Close()
}

// ServerOptions represents the tunable options this server provides
type ServerOptions struct {
	Pagerduty pagerduty.ServiceOptions
	// TODO - actually implement this lol
	Port int
}
