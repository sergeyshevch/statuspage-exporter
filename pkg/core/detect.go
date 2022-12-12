package core

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
	"github.com/sergeyshevch/statuspage-exporter/pkg/statusio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/statuspageio"
)

// DetectStatusPageType detects statuspage engine for all configured statuspage URLs.
func DetectStatusPageType(log *zap.Logger, restyClient *resty.Client) ([]string, []string) {
	targetUrls := config.StatusPages()

	var statusPageIoPages, statusIoPages []string

	for _, targetURL := range targetUrls {
		if statuspageio.IsStatusPageIOPage(log, targetURL, restyClient) {
			log.Info("Detected StatusPage.io page", zap.String("url", targetURL))
			statusPageIoPages = append(statusPageIoPages, targetURL)
		} else if statusio.IsStatusIOPage(log, targetURL, restyClient) {
			log.Info("Detected Status.io page", zap.String("url", targetURL))
			statusIoPages = append(statusIoPages, targetURL)
		}
	}

	return statusPageIoPages, statusIoPages
}
