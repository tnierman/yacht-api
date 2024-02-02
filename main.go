package main

import (
	"encoding/json"
	"log"
	"net/http"

	gopagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/gorilla/mux"
	"github.com/tnierman/yacht-api/pkg/pagerduty"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/clusters/{id}", clusterHandler)

	server := &http.Server {
		Handler: r,
		Addr: "127.0.0.1:8000",
	}
	log.Fatal(server.ListenAndServe())
}

// ClusterResponse defines the response body for a request at /api/clusters/{id}
type ClusterResponse struct {
	ClusterID string								 `json:"ID,omitempty"`
	Incidents []gopagerduty.Incident `json:"incidents,omitempty"`
}

func clusterHandler(resp http.ResponseWriter, req *http.Request) {
	body := ClusterResponse{}

	vars := mux.Vars(req)
	clusterID, found := vars["id"]
	if !found {
		// TODO - log and return error in resp
		log.Fatalf("failed to determine cluster ID from request: %#v", req)
	}
	body.ClusterID = clusterID

	incidents, err := pagerduty.GetClusterIncidents(clusterID)
	if err != nil {
		log.Fatalf("failed to retrieve incidents from PagerDuty: %v", err)
	}
	body.Incidents = incidents

	bytes, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("failed to encode response body: %v", err)
	}
	resp.Write(bytes)
}
