package statuspageio

import (
	"net/url"
	"sync"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/metrics"
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
		metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()
		log.Error("Failed to check if page is StatusPage.io page", zap.String("url", targetURL), zap.Error(err))
	}

	return resp.IsSuccess()
}

// FetchStatusPages fetch given status pages and export result as metrics.
func FetchStatusPages(log *zap.Logger, client *resty.Client, targetUrls []string) {
	wg := &sync.WaitGroup{}

	for _, targetURL := range targetUrls {
		go fetchStatusPage(wg, log, targetURL, client)
	}

	wg.Wait()
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

func fetchStatusPage(wg *sync.WaitGroup, log *zap.Logger, targetURL string, client *resty.Client) {
	wg.Add(1)
	log.Info("Fetching status page", zap.String("url", targetURL))

	defer wg.Done()

	parsedURL, err := constructURL(log, targetURL)
	if err != nil {
		metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()

		return
	}

	resp, err := client.R().SetResult(&AtlassianStatusPageResponse{}).Get(parsedURL) //nolint:exhaustruct
	if err != nil {
		log.Error(
			"Error fetching status page",
			zap.String("url", targetURL),
			zap.Duration("duration", resp.Time()),
			zap.Error(err),
		)
		metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()

		return
	}

	result, ok := resp.Result().(*AtlassianStatusPageResponse)
	if !ok {
		log.Error(
			"Error parsing status page response",
			zap.String("url", targetURL),
			zap.Duration("duration", resp.Time()),
			zap.Error(err),
		)
		metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()

		return
	}

	for _, component := range result.Components {
		metrics.ServiceStatus.WithLabelValues(
			result.Page.Name,
			targetURL,
			component.Name,
		).Set(utils.StatusToMetricValue(component.Status))
	}

	log.Info("Fetched status page", zap.Duration("duration", resp.Time()), zap.String("url", targetURL))
}
