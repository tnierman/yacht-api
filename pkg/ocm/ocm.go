package ocm

import (
	"fmt"

	"github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/osdctl/pkg/utils"
)

func NewClient() (*sdk.Connection, error) {
	return utils.CreateConnection()
}

func GetCluster(clusterID string) (*cmv1.Cluster, error) {
	client, err := NewClient()
	if err != nil {
		return &cmv1.Cluster{}, fmt.Errorf("failed to initialize client: %w", err)
	}
	return utils.GetCluster(client, clusterID)
}
