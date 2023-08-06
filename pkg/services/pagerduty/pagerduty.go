package pagerduty

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tnierman/yacht-api/pkg/logging"
	"github.com/tnierman/yacht-api/pkg/services/pagerduty/collector"
)

type Service struct {
	collector *collector.Collector
	logger *logging.Logger
}

func NewService(parent *logging.Logger, opts ServiceOptions) *Service {
	c := collector.NewCollector(opts.AuthToken, opts.Teams)
	l := parent.NewChildLogger("pagerduty")
	s := &Service{
		collector: c,
		logger: l,
	}
	return s
}

func (s *Service) AddRoutes(router *mux.Router) {
	router.HandleFunc("/incidents", s.GetIncidents)
}

func (s *Service) Start(logger *logging.Logger) {
	s.collector.Collect(logger)
}

func (s *Service) Stop() {
	s.collector.Stop()
}

func (s Service) GetIncidents(w http.ResponseWriter, _ *http.Request) {
	data, err := json.Marshal(s.collector.Incidents())
	if err != nil {
		http.Error(w, "could not decode data", http.StatusInternalServerError)
		s.logger.Sugar().Errorf("failed to decode incident data: %v", err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		s.logger.Sugar().Errorf("failed to write response to request: %v", err)
	}
}

type ServiceOptions struct {
	AuthToken string
	Teams     []string
}
