package statuspageio

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/utils"
)

// IsStatusPageIOPage checks if given URL is StatusPage.io page.
func IsStatusPageIOPage(log *zap.Logger, targetURL string, client *resty.Client) bool {
	parsedURL, err := constructURL(log, targetURL)
	if err != nil {
		return false
	}

	resp, err := client.R().Head(parsedURL)
	if err != nil {
		log.Error(
			"Failed to check if page is StatusPage.io page",
			zap.String("url", targetURL),
			zap.Error(err),
		)
	}

	return resp.IsSuccess()
}

func constructURL(log *zap.Logger, targetURL string) (string, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	parsedURL.Path = "/api/v2/components.json"

	if parsedURL.Host == "" {
		log.Error(utils.ErrInvalidURL.Error(), zap.String("url", targetURL))

		return "", utils.ErrInvalidURL
	}

	return parsedURL.String(), nil
}

// FetchStatusPage fetches status page and updates service status gauge.
func FetchStatusPage(
	log *zap.Logger,
	targetURL string,
	client *resty.Client,
	componentStatus *prometheus.GaugeVec,
	overallStatus *prometheus.GaugeVec,
) error {
	log.Info("Fetching status page", zap.String("url", targetURL))

	parsedURL, err := constructURL(log, targetURL)
	if err != nil {
		return err
	}

	resp, err := client.R().
		SetResult(&AtlassianStatusPageResponse{}). //nolint:exhaustruct
		Get(parsedURL)
	if err != nil {
		log.Error(
			"Error fetching status page",
			zap.String("url", targetURL),
			zap.Duration("duration", resp.Time()),
			zap.Error(err),
		)

		return err
	}

	result, ok := resp.Result().(*AtlassianStatusPageResponse)
	if !ok {
		log.Error(
			"Error parsing status page response",
			zap.String("url", targetURL),
			zap.Duration("duration", resp.Time()),
			zap.Error(err),
		)

		return err
	}

	for _, component := range result.Components {
		componentStatus.WithLabelValues(
			result.Page.Name,
			targetURL,
			component.Name,
		).Set(float64(IndicatorToMetricValue(component.Status)))
	}

	overallStatus.WithLabelValues(
		result.Page.Name,
		targetURL,
	).Set(float64(IndicatorToMetricValue(result.Status.Indicator)))

	log.Info(
		"Fetched status page",
		zap.Duration("duration", resp.Time()),
		zap.String("url", targetURL),
	)

	return nil
}
