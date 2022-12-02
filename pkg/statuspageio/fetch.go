package statuspageio

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
	"github.com/sergeyshevch/statuspage-exporter/pkg/metrics"
	"github.com/sergeyshevch/statuspage-exporter/pkg/utils"
)

// StartFetchingLoop starts a loop that fetches status pages.
func StartFetchingLoop(ctx context.Context, wg *sync.WaitGroup, log *zap.Logger) {
	wg.Add(1)
	defer wg.Done()

	fetchDelay := config.FetchDelay()
	client := resty.New().EnableTrace().SetTimeout(config.ClientTimeout())

	for {
		select {
		default:
			fetchAllStatusPages(log, client)

			time.Sleep(fetchDelay)
		case <-ctx.Done():
			log.Info("Stopping fetching loop")

			return
		}
	}
}

func fetchAllStatusPages(log *zap.Logger, client *resty.Client) {
	wg := &sync.WaitGroup{}

	targetUrls := config.StatusPageIoPages()

	for _, targetURL := range targetUrls {
		go fetchStatusPage(wg, log, targetURL, client)
	}

	wg.Wait()
}

func fetchStatusPage(wg *sync.WaitGroup, log *zap.Logger, targetURL string, client *resty.Client) {
	wg.Add(1)
	log.Info("Fetching status page", zap.String("url", targetURL))

	defer wg.Done()

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		panic(err)
	}

	parsedURL.Path = "/api/v2/components.json"

	if parsedURL.Host == "" {
		log.Error("Invalid URL. It won't be parsed. Check that your url contains scheme", zap.String("url", targetURL))
		metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()

		return
	}

	resp, err := client.R().SetResult(&AtlassianStatusPageResponse{}).Get(parsedURL.String()) //nolint:exhaustruct
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
