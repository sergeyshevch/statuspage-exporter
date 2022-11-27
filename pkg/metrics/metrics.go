package metrics

import "github.com/prometheus/client_golang/prometheus"

var ServiceStatusFetchError = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "service_status_fetch_error",
		Help: "Number of errors encountered while fetching service status",
	},
	[]string{"url"},
)

var ServiceStatus = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "service_status",
		Help: "Status of a service component, values 0 (operational) to 4 (major_outage)",
	},
	[]string{"service", "status_page_url", "component"},
)
