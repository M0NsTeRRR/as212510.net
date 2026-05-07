package main

import (
	"log/slog"

	"github.com/m0nsterrr/as212510.net/v3/internal/config"
	"github.com/m0nsterrr/as212510.net/v3/internal/logging"
	"github.com/m0nsterrr/as212510.net/v3/internal/server"
)

var (
	version   = "development"
	buildTime = "0"
)

func main() {
	logging.Init()

	slog.Info("Starting as212510", "version", version, "build_time", buildTime)

	config := config.Init()

	go server.StartHealthcheckServer(config.HealthCheck.Address)
	go server.StartMetricServer(config.Metric.Address)
	server.StartServer(config)
}
