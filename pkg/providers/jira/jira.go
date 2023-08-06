package jira

import "github.com/gorilla/mux"

type Service struct {
}

func NewService() Service {
	s := Service{}
	return s
}

func (s *Service) AddRoutes(router *mux.Router) {
}
