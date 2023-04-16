package statusio

import "github.com/sergeyshevch/statuspage-exporter/pkg/engines/types"

// StatusToMetricValue converts statuspage status to metric value.
func StatusToMetricValue(status string) types.Status {
	switch status {
	case "Operational":
		return types.OperationalStatus
	case "Planned Maintenance":
		return types.PlannedMaintenanceStatus
	case "Degraded Performance":
		return types.DegradedPerformanceStatus
	case "Partial Service Disruption":
		return types.PartialOutageStatus
	case "Service Disruption":
		return types.MajorOutageStatus
	case "Security Issue":
		return types.SecurityIssueStatus
	default:
		return types.UnknownStatus
	}
}

// PageDescriptionToMetricValue converts statuspage status to metric value.
func PageDescriptionToMetricValue(description string) types.Status {
	switch description {
	case "All Systems Operational":
		return types.OperationalStatus
	case "Planned Maintenance In Progress":
		return types.PlannedMaintenanceStatus
	case "Active Issue":
		return types.PartialOutageStatus
	default:
		return types.UnknownStatus
	}
}
