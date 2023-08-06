package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tnierman/yacht-api/pkg/apiserver"
	"github.com/tnierman/yacht-api/pkg/logging"
	"github.com/tnierman/yacht-api/pkg/services/pagerduty"
)

func main() {
	logger, err := logging.NewLogger("main")
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Cleanup()

	pdAuthToken, err := getPDAuthToken()
	if err != nil {
		logger.Sugar().Fatalf("failed to retrieve token from ~/.pd.yml: %v", err)
	}

	opts := apiserver.ServerOptions{
		Pagerduty: pagerduty.ServiceOptions{
			AuthToken: pdAuthToken,
			Teams: []string{"Platform SRE"},
		},
	}
	server := apiserver.NewServer(logger, opts)
	err = server.Serve(logger)
	if err != nil {
		logger.Sugar().Fatalf("failed to run apiserver: %v", err)
	}
}

func getPDAuthToken() (string, error) {
	pdToken, found := os.LookupEnv("PD_TOKEN")
	if !found || pdToken == "" {
		return "", fmt.Errorf("$PD_TOKEN unset")
	}
	return pdToken, nil
}
