package jira

import (
	"fmt"
	"os"

	jira "github.com/andygrunwald/go-jira"
	"github.com/tnierman/yacht-api/pkg/ocm"
)

const (
	jiraTokenEnvVar string = "JIRA_TOKEN"
	JiraBaseURL string = "https://issues.redhat.com"
)

func NewClient() (*jira.Client, error) {
	jiraToken := os.Getenv(jiraTokenEnvVar)
	if jiraToken == "" {
		return nil, fmt.Errorf("'%s' is unset", jiraTokenEnvVar)
	}
	transport := jira.PATAuthTransport {
		Token: jiraToken,
	}
	return jira.NewClient(transport.Client(), JiraBaseURL)
}

func GetClusterTickets(clusterID string) ([]jira.Issue, error) {
	client, err := NewClient()
	if err != nil {
		return []jira.Issue{}, fmt.Errorf("failed to initialize new client: %w", err)
	}

	// TODO: cache eventually
	cluster, err := ocm.GetCluster(clusterID)
	if err != nil {
		return []jira.Issue{}, fmt.Errorf("failed to retrieve cluster from OCM: %w", err)
	}

	externalClusterID := cluster.ExternalID()

	fmt.Printf("externalClusterID: %v\n", externalClusterID)
	fmt.Printf("clusterID: %v\n", clusterID)

	jql := fmt.Sprintf(
		`(project = "OpenShift Hosted SRE Support" AND "Cluster ID" ~ "%s") OR (project = "OpenShift Hosted SRE Support" AND "Cluster ID" ~ "%s") ORDER BY created DESC`,
		externalClusterID,
		clusterID,
	)

	issues, response, err := client.Issue.Search(jql, &jira.SearchOptions{})
	fmt.Printf("issues: %v\n", issues)
	fmt.Printf("response: %v\n", response)
	if err != nil {
		return []jira.Issue{}, fmt.Errorf("failed to retrieve issues from %s: %w", JiraBaseURL, err)
	}
	return issues, nil
}
