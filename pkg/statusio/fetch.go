package statusio

import (
	"context"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/PuerkitoBio/goquery"

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

	targetUrls := config.StatusIoPages()

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

	parsedURL.Path = "/"

	if parsedURL.Host == "" {
		log.Error("Invalid URL. It won't be parsed. Check that your url contains scheme", zap.String("url", targetURL))
		metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()

		return
	}

	resp, err := client.R().SetDoNotParseResponse(true).Get(parsedURL.String())
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

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Error("Error closing response body", zap.Error(err))
			metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()

			return
		}
	}(resp.RawBody())

	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		log.Error("Error parsing response body", zap.Error(err))
		metrics.ServiceStatusFetchError.WithLabelValues(targetURL).Inc()

		return
	}

	title := strings.Split(doc.Find("title").Text(), " ")[0]

	doc.Find(".component").Each(func(_ int, s *goquery.Selection) {
		componentName := s.Find(".component_name").First().Text()
		componentStatus := s.Find(".component-status").First().Text()

		metrics.ServiceStatus.WithLabelValues(
			title,
			targetURL,
			strings.Trim(componentName, " "),
		).Set(utils.StatusToMetricValue(componentStatus))
	})

	log.Info("Fetched status page", zap.Duration("duration", resp.Time()), zap.String("url", targetURL))
}
