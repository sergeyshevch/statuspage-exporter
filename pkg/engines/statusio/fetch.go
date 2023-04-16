package statusio

import (
  "io"
  "net/url"
  "strings"

  "github.com/go-resty/resty/v2"
  "github.com/prometheus/client_golang/prometheus"
  "go.uber.org/zap"

  "github.com/PuerkitoBio/goquery"

  "github.com/sergeyshevch/statuspage-exporter/pkg/utils"
)

// IsStatusIOPage checks if given URL is Status.io page.
func IsStatusIOPage(log *zap.Logger, targetURL string, client *resty.Client) bool {
  parsedURL, err := constructURL(log, targetURL)
  if err != nil {
    return false
  }

  resp, err := client.R().Get(parsedURL)
  if err != nil {
    log.Error("Failed to check if page is Status.io page", zap.String("url", targetURL), zap.Error(err))

  }

  return strings.Contains(resp.String(), "status.io")
}

func constructURL(log *zap.Logger, targetURL string) (string, error) {
  parsedURL, err := url.Parse(targetURL)
  if err != nil {
    panic(err)
  }

  parsedURL.Path = "/"

  if parsedURL.Host == "" {
    log.Error(utils.ErrInvalidURL.Error(), zap.String("url", targetURL))

    return "", utils.ErrInvalidURL
  }

  return parsedURL.String(), nil
}

func FetchStatusPage(log *zap.Logger, targetURL string, client *resty.Client, serviceStatusGauge *prometheus.GaugeVec) error {
  log.Info("Fetching status page", zap.String("url", targetURL))

  parsedURL, err := constructURL(log, targetURL)
  if err != nil {
    return err
  }

  resp, err := client.R().SetDoNotParseResponse(true).Get(parsedURL)
  if err != nil {
    log.Error(
      "Error fetching status page",
      zap.String("url", targetURL),
      zap.Duration("duration", resp.Time()),
      zap.Error(err),
    )

    return err
  }

  defer func(body io.ReadCloser) {
    err := body.Close()
    if err != nil {
      log.Error("Error closing response body", zap.Error(err))

      return
    }
  }(resp.RawBody())

  doc, err := goquery.NewDocumentFromReader(resp.RawBody())
  if err != nil {
    log.Error("Error parsing response body", zap.Error(err))

    return err
  }

  title := strings.Split(doc.Find("title").Text(), " ")[0]

  doc.Find(".component").Each(func(_ int, s *goquery.Selection) {
    componentName := s.Find(".component_name").First().Text()
    componentStatus := s.Find(".component-status").First().Text()

    serviceStatusGauge.WithLabelValues(
      title,
      targetURL,
      strings.Trim(componentName, " "),
    ).Set(utils.StatusToMetricValue(componentStatus))
  })

  log.Info("Fetched status page", zap.Duration("duration", resp.Time()), zap.String("url", targetURL))

  return nil
}
