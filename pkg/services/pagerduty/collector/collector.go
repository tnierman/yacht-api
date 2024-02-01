package collector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/tnierman/yacht-api/pkg/logging"
)

type Collector struct {
	// Internals
	client *pagerduty.Client

	ctx context.Context
	context.CancelFunc
	stopped sync.WaitGroup

	// Collection settings
	Teams []string

	// Collection objects
	incidents []pagerduty.Incident
	incidentCollection sync.WaitGroup
}

func NewCollector(authToken string, teams []string) *Collector {
	// Initialize internals
	client := pagerduty.NewClient(authToken)
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize collection objects
	incidents := []pagerduty.Incident{}

	c := &Collector{
		// Internals
		client: client,
		ctx: ctx,
		CancelFunc: cancel,
		stopped: sync.WaitGroup{},

		// Collection settings
		Teams: teams,

		// Collection objects
		incidents: incidents,
	}
	return c
}

func (c *Collector) Collect(logger *logging.Logger) {
	go c.collect(logger)
}

func (c *Collector) collect(logger *logging.Logger) {
	c.stopped.Add(1)
	for {
		// Check if we need to stop
		select {
		case <- c.ctx.Done():
			c.stopped.Done()
			return
		default:
		}

		// Collect objects
		err := c.collectIncidents()
		if err != nil {
			logger.Sugar().Errorf("failed to collect one or more incidents: %v", err)
		}

		time.Sleep(5*time.Second)
	}
}

// collectIncidents refreshes the Collector's internal cache with the latest active, high-urgency incidents from the Collector's team
func (c *Collector) collectIncidents() error {
	c.incidentCollection.Wait()
	c.incidentCollection.Add(1)
	// Clear incidents and repopulate
	c.incidents = []pagerduty.Incident{}
	err := c.collectIncidentsFromOffset(0)
	c.incidentCollection.Done()
	return err
}

// collectIncidentsFromOffset gathers all remaining incidents starting from the offset
func (c *Collector) collectIncidentsFromOffset(offset int) error {
	response, err := c.client.ListIncidentsWithContext(c.ctx, pagerduty.ListIncidentsOptions{
		TeamIDs: c.Teams,
		Statuses: []string{"acknowledged"},
		Urgencies: []string{"high"},
		Offset: uint(offset),
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve incidents from pagerduty: %w", err)
	}

	c.incidents = append(c.incidents, response.Incidents...)
	if response.More {
		return c.collectIncidentsFromOffset(len(c.incidents))
	}
	return nil
}

func (c *Collector) Incidents() []pagerduty.Incident {
	c.incidentCollection.Wait()
	c.incidentCollection.Add(1)
	incidents := c.incidents
	c.incidentCollection.Done()
	return incidents
}

// Stop blocks until the collector has fully stopped collecting metrics
func (c *Collector) Stop() {
	c.CancelFunc()
	c.stopped.Wait()
}

