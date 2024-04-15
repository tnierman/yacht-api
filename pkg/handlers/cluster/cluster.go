package cluster

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tnierman/yacht-api/pkg/jira"
	"github.com/tnierman/yacht-api/pkg/pagerduty"
	clusterapi "github.com/tnierman/yacht-api/api/cluster"
)

func ClusterIDHandler(resp http.ResponseWriter, req *http.Request) {
	response := clusterapi.Response{}

	vars := mux.Vars(req)
	clusterID, found := vars["id"]
	if !found {
		log.Fatalf("failed to determine cluster ID from request: %#v", req)
	}
	response.ClusterID = clusterID

	// handle PagerDuty
	incidents, err := pagerduty.GetClusterIncidents(clusterID)
	if err != nil {
		log.Fatalf("failed to retrieve incidents from PagerDuty: %v", err)
	}
	response.PagerDuty.Incidents = incidents

	// handle Jira
	issues, err := jira.GetClusterTickets(clusterID)
	if err != nil {
		log.Fatalf("failed to retrieve jira tickets: %v", err)
	}

	fmt.Printf("issues: %v\n", issues)
	// TODO: convert
	//response.Jira.Issues = issues

	bytes, err := response.Bytes()
	if err != nil {
		log.Fatalf("failed to encode response body: %v", err)
	}
	resp.Write(bytes)
}
