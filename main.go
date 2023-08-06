package main

import (
	"log"

	"github.com/tnierman/yacht-api/pkg/api"
	"github.com/tnierman/yacht-api/pkg/logging"
)

func main() {
	logger, err := logging.NewLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Cleanup()

	server := api.NewServer()
	err = server.Serve(logger)
	if err != nil {
		
	}
}
