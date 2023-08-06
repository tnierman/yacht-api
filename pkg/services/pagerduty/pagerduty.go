package pagerduty

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Service struct {

}

func NewService() Service {
	s := Service{}
	return s
}

func (s *Service) AddRoutes(router *mux.Router) {
	router.HandleFunc("/incidents", s.GetIncidents)
}

func (s Service) GetIncidents(w http.ResponseWriter, r *http.Request) {
}
