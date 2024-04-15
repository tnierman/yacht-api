package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/gorilla/mux"
	"github.com/tnierman/yacht-api/pkg/ocm"
)

// NewClient initializes a PagerDuty client using the API token specified with $PAGERDUTY_TOKEN
func NewClient() (*pagerduty.Client, error) {
	token, found := os.LookupEnv("PAGERDUTY_TOKEN")
	if !found {
		return &pagerduty.Client{}, fmt.Errorf("$PAGERDUTY_TOKEN unset")
	}
	return pagerduty.NewClient(token), nil
}

// GetClusterIncidents fetches all acknowledged, high-urgency (paging) alerts for a cluster
func GetClusterIncidents(clusterID string) ([]pagerduty.Incident, error) {
	client, err := NewClient()
	if err != nil {
		return []pagerduty.Incident{}, fmt.Errorf("failed to initialize PagerDuty client: %w", err)
	}
	service, err := GetServiceForCluster(clusterID, client)
	if err != nil {
		return []pagerduty.Incident{}, fmt.Errorf("failed to retrieve PagerDuty service for cluster '%s': %w", clusterID, err)
	}
	incidents, err := GetIncidentsForService(service.ID, client)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch incidents from PagerDuty: %w", err)
	}
	return incidents, nil
}

// GetServiceForCluster fetches the PagerDuty service for a cluster
func GetServiceForCluster(clusterID string, client *pagerduty.Client) (pagerduty.Service, error) {
	cluster, err := ocm.GetCluster(clusterID)
	if err != nil {
		return pagerduty.Service{}, fmt.Errorf("failed to retrieve cluster from OCM: %w", err)
	}
	query, found := cluster.DNS().GetBaseDomain()
	if !found {
		return pagerduty.Service{}, fmt.Errorf("failed to retrieve base domain for cluster '%s'. Cluster details: %#v", clusterID, cluster)
	}
	response, err := client.ListServicesWithContext(context.TODO(), pagerduty.ListServiceOptions{
		// TODO - figure out a way to pass in SRE-P team ID
		// TeamIDs: []string{""},
		Query: cluster.DNS().BaseDomain(),
	})
	if err != nil {
		return pagerduty.Service{}, fmt.Errorf("failed to list services matching '%s': %w", query, err)
	}

	results := len(response.Services)
	if results != 1 {
		return pagerduty.Service{}, fmt.Errorf("failed to locate service matching '%s': expected 1 result, got %d", query, results)
	}
	return response.Services[0], nil
}

// GetIncidentsForService fetches all acknowledged, high-urgency (paging) alerts for a PagerDuty service
func GetIncidentsForService(serviceID string, client *pagerduty.Client) ([]pagerduty.Incident, error) {
	return getIncidentsForServiceAtOffset(serviceID, client, 0)
}

// getIncidentsForServiceAtOffset fetches acknowledged, high-urgency (paging) alerts for a PagerDuty service starting at the provided offset
func getIncidentsForServiceAtOffset(serviceID string, client *pagerduty.Client, offset uint) ([]pagerduty.Incident, error) {
	fmt.Println("serviceID: ", serviceID)
	response, err := client.ListIncidentsWithContext(context.TODO(), pagerduty.ListIncidentsOptions{
		ServiceIDs: []string{serviceID},
		// TODO: expose these values in an option-set (ie - create a GetIncidentsOptions struct)
		Urgencies: []string{"high"},
		Statuses: []string{"acknowledged"},
		//Statuses: []string{"resolved"},
		Offset: offset,
	})
	if err != nil {
		return []pagerduty.Incident{}, fmt.Errorf("failed to retrieve incidents for PagerDuty service '%s': %w", serviceID, err)
	}
	fmt.Printf("response: %#v\n", response)

	incidents := response.Incidents
	if response.More {
		offset += uint(len(response.Incidents))
		more, err := getIncidentsForServiceAtOffset(serviceID, client, offset)
		if err != nil {
			return []pagerduty.Incident{}, fmt.Errorf("failed to retrieve additional incidents starting at offset '%d': %w", offset, err)
		}
		incidents = append(incidents, more...)
	}
	return incidents, nil
}

type ClusterResponse struct {
	ClusterID string
	Incidents []pagerduty.Incident
}

func ClusterIDHandler(resp http.ResponseWriter, req *http.Request) {
	body := ClusterResponse{}

	vars := mux.Vars(req)
	clusterID, found := vars["id"]
	if !found {
		// TODO - log and return error in resp
		log.Fatalf("failed to determine cluster ID from request: %#v", req)
	}
	body.ClusterID = clusterID

	incidents, err := GetClusterIncidents(clusterID)
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
