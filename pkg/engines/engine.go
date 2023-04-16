package engines

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statusio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/statuspageio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/engines/types"
)

var errUnknownStatusPageType = fmt.Errorf("unknown statuspage type")

// FetchStatus detect statuspage type and fetch its status.
func FetchStatus(
	log *zap.Logger,
	targetURL string,
	componentStatus *prometheus.GaugeVec,
	overallStatus *prometheus.GaugeVec,
) error {
	restyClient := resty.New().
		EnableTrace().
		SetTimeout(config.ClientTimeout()).
		SetRetryCount(config.RetryCount())

	statusPageType := DetectStatusPageType(log, restyClient, targetURL)
	if statusPageType == types.UnknownType {
		return errUnknownStatusPageType
	}

	switch statusPageType {
	case types.StatusPageIOType:
		return statuspageio.FetchStatusPage(
			log,
			targetURL,
			restyClient,
			componentStatus,
			overallStatus,
		)
	case types.StatusIOType:
		return statusio.FetchStatusPage(log, targetURL, restyClient, componentStatus, overallStatus)
	case types.UnknownType:
		return errUnknownStatusPageType
	default:
		return errUnknownStatusPageType
	}
}
