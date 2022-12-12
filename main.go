package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
	"github.com/sergeyshevch/statuspage-exporter/pkg/core"
	"github.com/sergeyshevch/statuspage-exporter/pkg/metrics"
	"github.com/sergeyshevch/statuspage-exporter/pkg/statusio"
	"github.com/sergeyshevch/statuspage-exporter/pkg/statuspageio"
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 5 * time.Second
)

func startHTTP(ctx context.Context, wg *sync.WaitGroup, log *zap.Logger) {
	wg.Add(1)
	defer wg.Done()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpPort := config.HTTPPort()
	httpServer := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf(":%d", httpPort),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           mux,
	}

	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic("Unable to start a http server.", zap.Error(err))
		}
	}()

	log.Info("Http server listening on", zap.Int("port", httpPort))

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil { //nolint:contextcheck
		log.Panic("Http server Shutdown Failed", zap.Error(err))
	}

	log.Info("Http server stopped")
}

func main() {
	log, err := config.InitConfig()
	if err != nil {
		log.Fatal("Unable to initialize config", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	prometheus.MustRegister(version.NewCollector("statuspage_exporter"))
	prometheus.MustRegister(metrics.ServiceStatus)
	prometheus.MustRegister(metrics.ServiceStatusFetchError)

	restyClient := resty.New().EnableTrace().SetTimeout(config.ClientTimeout()).SetRetryCount(config.RetryCount())

	statusPageIOTargets, statusIOTargets := core.DetectStatusPageType(log, restyClient)

	go core.StartFetchingLoop(ctx, wg, log, restyClient, statusPageIOTargets, statuspageio.FetchStatusPages)
	go core.StartFetchingLoop(ctx, wg, log, restyClient, statusIOTargets, statusio.FetchStatusPages)
	go startHTTP(ctx, wg, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("Received shutdown signal. Waiting for workers to terminate...")
	cancel()

	wg.Wait()
}
