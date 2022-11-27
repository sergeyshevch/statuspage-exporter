package metrics

import "github.com/prometheus/client_golang/prometheus"

// ServiceStatusFetchError is a counter that counts errors while fetching service status.
var ServiceStatusFetchError = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        "fetch_error_total",
		Namespace:   "service",
		Subsystem:   "status",
		Help:        "Number of errors encountered while fetching service status",
		ConstLabels: map[string]string{},
	},
	[]string{"url"},
)

// ServiceStatus is a gauge that represents service status.
var ServiceStatus = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name:        "",
		Namespace:   "service",
		Subsystem:   "status",
		Help:        "Status of a service component, values 0 (operational) to 4 (major_outage)",
		ConstLabels: map[string]string{},
	},
	[]string{"service", "status_page_url", "component"},
)
