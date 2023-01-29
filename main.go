package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	cpuTemp        prometheus.Gauge
	sbi3gppHeaders *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),
		sbi3gppHeaders: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "sbi_3gpp_headers_attempt",
				Help: "Number of times request sent using 3gpp headers.",
			},
			[]string{"header_name"},
		),
	}
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.sbi3gppHeaders)
	return m
}

func main() {
	// http.Handle("/metrics", promhttp.Handler())
	// panic(http.ListenAndServe(":8080", nil))

	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)
	// Set values for the new created metrics.
	m.cpuTemp.Set(65.3)

	m.sbi3gppHeaders.With(prometheus.Labels{"header_name": "target_api_root"}).Inc()
	m.sbi3gppHeaders.With(prometheus.Labels{"header_name": "routing_binding"}).Inc()

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
