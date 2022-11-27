package statuspage

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
	"github.com/sergeyshevch/statuspage-exporter/pkg/metrics"
)

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

	targetUrls := config.StatusPages()

	for _, targetUrl := range targetUrls {
		go fetchStatusPage(wg, log, targetUrl, client)
	}
	wg.Wait()
}

func fetchStatusPage(wg *sync.WaitGroup, log *zap.Logger, targetUrl string, client *resty.Client) {
	wg.Add(1)
	log.Info("Fetching status page", zap.String("url", targetUrl))
	defer wg.Done()

	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		panic(err)
	}
	parsedUrl.Path = "/api/v2/components.json"
	if parsedUrl.Host == "" {
		log.Error("Invalid URL. It won't be parsed. Check that your url contains scheme", zap.String("url", targetUrl))
		metrics.ServiceStatusFetchError.WithLabelValues(targetUrl).Inc()
		return
	}

	resp, err := client.R().SetResult(&AtlassianStatusPageResponse{}).Get(parsedUrl.String())
	if err != nil {
		log.Error("Error fetching status page", zap.String("url", targetUrl), zap.Duration("duration", resp.Time()), zap.Error(err))
		metrics.ServiceStatusFetchError.WithLabelValues(targetUrl).Inc()
		return
	}

	result := resp.Result().(*AtlassianStatusPageResponse)
	for _, component := range result.Components {
		metrics.ServiceStatus.WithLabelValues(result.Page.Name, targetUrl, component.Name).Set(statusToMetricValue(component.Status))
	}

	log.Info("Fetched status page", zap.Duration("duration", resp.Time()), zap.String("url", targetUrl))
}

func statusToMetricValue(status string) float64 {
	switch status {
	case "operational":
		return 1
	case "degraded_performance":
		return 2
	case "partial_outage":
		return 3
	case "major_outage":
		return 4
	default:
		return 0
	}
}
