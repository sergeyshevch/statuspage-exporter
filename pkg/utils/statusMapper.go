package utils

// StatusToMetricValue converts statuspage status to metric value.
func StatusToMetricValue(status string) float64 {
	switch status {
	case "operational":
		return 1
	case "degraded_performance":
		return 2 //nolint:gomnd
	case "partial_outage":
		return 3 //nolint:gomnd
	case "major_outage":
		return 4 //nolint:gomnd
	default:
		return 0
	}
}
