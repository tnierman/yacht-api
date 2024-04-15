package pagerduty

import (
	gopagerduty "github.com/PagerDuty/go-pagerduty"
)

type Response struct {
	Service   string                 `json:"service"`
	Incidents []gopagerduty.Incident `json:"incidents"`
}
