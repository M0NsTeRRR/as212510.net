package main

import (
	"log"
	"net/http"
)

var (
	healthRouter = http.NewServeMux()
)

func startHealthcheck(address string) {
	healthRouter.HandleFunc("/healthcheck", healtcheck)
	log.Printf("Healthcheck handler is listening on %s", address)
	log.Fatal(http.ListenAndServe(address, healthRouter))
}

func healtcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
