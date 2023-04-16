package engines

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statusio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statuspageio"
)

var errUnknownStatusPageType = fmt.Errorf("unknown statuspage type")

// FetchStatus detect statuspage type and fetch its status.
func FetchStatus(log *zap.Logger, targetURL string, serviceStatusGauge *prometheus.GaugeVec) error {
	restyClient := resty.New().
		EnableTrace().
		SetTimeout(config.ClientTimeout()).
		SetRetryCount(config.RetryCount())

	statusPageType := DetectStatusPageType(log, restyClient, targetURL)
	if statusPageType == Unknown {
		return errUnknownStatusPageType
	}

	switch statusPageType {
	case StatusPageIO:
		return statuspageio.FetchStatusPage(log, targetURL, restyClient, serviceStatusGauge)
	case StatusIO:
		return statusio.FetchStatusPage(log, targetURL, restyClient, serviceStatusGauge)
	case Unknown:
		return errUnknownStatusPageType
	default:
		return errUnknownStatusPageType
	}
}
