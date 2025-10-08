package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/halkeye/bambulabs-exporter/internal/exporter"
)

func main() {
	// Unregister default collectors
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.Unregister(collectors.NewGoCollector())

	// Create and start the exporter
	exp := exporter.NewExporter()
	
	// Connect to MQTT broker
	exp.ConnectToBroker()
	
	// Start HTTP server
	exp.StartHTTPServer()
	
	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":9101", nil))
}