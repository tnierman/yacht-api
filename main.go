package main

import (
	"log"

	"github.com/tnierman/yacht-api/pkg/apiserver"
	"github.com/tnierman/yacht-api/pkg/logging"
)

func main() {
	logger, err := logging.NewLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Cleanup()

	server := apiserver.NewServer()
	err = server.Serve(logger)
	if err != nil {
		
	}
}
