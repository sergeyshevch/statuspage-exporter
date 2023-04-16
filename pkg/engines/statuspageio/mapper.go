package statuspageio

import "github.com/sergeyshevch/statuspage-exporter/pkg/engines/types"

// StatusToMetricValue converts statuspage status to metric value.
func StatusToMetricValue(status string) types.Status {
	switch status {
	case "operational":
		return types.OperationalStatus
	case "degraded_performance":
		return types.DegradedPerformanceStatus
	case "partial_outage":
		return types.PartialOutageStatus
	case "major_outage":
		return types.MajorOutageStatus
	default:
		return types.UnknownStatus
	}
}

// IndicatorToMetricValue converts statuspage indicator to metric value.
func IndicatorToMetricValue(indicator string) types.Status {
	switch indicator {
	case "none":
		return types.OperationalStatus
	case "minor":
		return types.DegradedPerformanceStatus
	case "major":
		return types.PartialOutageStatus
	case "critical":
		return types.MajorOutageStatus
	default:
		return types.UnknownStatus
	}
}
