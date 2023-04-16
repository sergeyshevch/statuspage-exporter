package prober

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/engines"
)

// Handler returns a http handler for /probe endpoint.
func Handler(log *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		targetURL := ctx.QueryParam("target")
		if targetURL == "" {
			return ctx.String(http.StatusBadRequest, "target is required")
		}

		start := time.Now()

		serviceStatus := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{ //nolint:exhaustruct
				Name: "service_status",
				Help: "Status of a service component, values 0 (operational) to 4 (major_outage)",
			},
			[]string{"service", "status_page_url", "component"},
		)

		serviceStatusDurationGauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{ //nolint:exhaustruct
				Name: "service_status_fetch_duration_seconds",
				Help: "Returns how long the service status fetch took to complete in seconds",
			},
			[]string{"status_page_url"},
		)

		registry := prometheus.NewRegistry()
		registry.MustRegister(serviceStatus)
		registry.MustRegister(serviceStatusDurationGauge)

		err := engines.FetchStatus(log, targetURL, serviceStatus)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		duration := time.Since(start).Seconds()

		serviceStatusDurationGauge.WithLabelValues(targetURL).Set(duration)

		h := echo.WrapHandler(
			promhttp.HandlerFor(registry, promhttp.HandlerOpts{}), //nolint:exhaustruct
		)

		return h(ctx)
	}
}
