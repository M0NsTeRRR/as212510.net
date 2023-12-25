package main

import (
	"log"
	"net/http"
)

var (
	healthRouter   = http.NewServeMux()
	healtcheckPort = ":10240"
)

func startHealthcheck() {
	healthRouter.HandleFunc("/healthcheck", healtcheck)
	log.Printf("Healthcheck handler is listening on ", healtcheckPort)
	log.Fatal(http.ListenAndServe(healtcheckPort, healthRouter))
}

func healtcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
