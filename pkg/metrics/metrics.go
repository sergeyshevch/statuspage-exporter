package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "statuspage"
)

type metricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

func newMetric(metricName string, docString string, t prometheus.ValueType, labelNames []string) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", metricName),
			docString,
			labelNames,
			nil,
		),
		Type: t,
	}
}

var ServiceStatusMetric = newMetric(
	"service_status",
	"Status of a service component, values 0 (operational) to 4 (major_outage)",
	prometheus.GaugeValue,
	[]string{"service", "status_page_url", "component"},
)

// ServiceStatusFetchErrorMetric is a counter that counts errors while fetching service status.
var ServiceStatusFetchErrorMetric = newMetric(
	"service_status_fetch_errors_total",
	"Number of errors while fetching service status",
	prometheus.CounterValue,
	[]string{"url"},
)
