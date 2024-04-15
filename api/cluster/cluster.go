package cluster

import (
	"encoding/json"

	"github.com/tnierman/yacht-api/api/jira"
	"github.com/tnierman/yacht-api/api/pagerduty"
)

type Response struct {
	ClusterID string             `json:"id"`
	PagerDuty pagerduty.Response `json:"pagerduty"`
	Jira      jira.Response      `json:"jira"`
}

func (r *Response) Bytes() ([]byte, error) {
	return json.Marshal(r)
}
