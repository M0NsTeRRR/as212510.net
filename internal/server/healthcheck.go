package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartHealthcheckServer(address string) {
	http.HandleFunc("/health", healtcheck)
	log.Printf("Healthcheck Server is listening on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func healtcheck(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "ok")
	if err != nil {
		log.Fatal(err.Error())
	}
}
