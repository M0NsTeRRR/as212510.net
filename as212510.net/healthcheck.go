package main

import (
	"log"
	"net/http"
)

func startHealthcheck(address string) {
	http.HandleFunc("/healthcheck", healtcheck)
	log.Printf("Healthcheck handler is listening on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func healtcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
