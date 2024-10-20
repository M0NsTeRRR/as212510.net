package as212510

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startMetric(address string) {
	log.Printf("Metric handler is listening on %s", address)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(address, nil))
}
