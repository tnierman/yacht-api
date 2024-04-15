package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tnierman/yacht-api/pkg/handlers/cluster"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/clusters/{id}", cluster.ClusterIDHandler)

	server := &http.Server {
		Handler: r,
		Addr: "127.0.0.1:8000",
	}
	log.Fatal(server.ListenAndServe())
}
