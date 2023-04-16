package main

import (
  "context"
  "fmt"
  "net/http"
  "os"
  "os/signal"
  "sync"
  "syscall"
  "time"

  echoPrometheus "github.com/labstack/echo-contrib/prometheus"
  "github.com/labstack/echo/v4"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/common/version"
  "go.uber.org/zap"

  "github.com/sergeyshevch/statuspage-exporter/pkg/config"
  "github.com/sergeyshevch/statuspage-exporter/pkg/prober"
)

const (
  shutdownTimeout = 5 * time.Second
)

func handleHealthz(ctx echo.Context) error {
  return ctx.String(http.StatusOK, "ok")
}

func startHTTP(ctx context.Context, wg *sync.WaitGroup, log *zap.Logger) {
  wg.Add(1)
  defer wg.Done()

  srv := echo.New()
  echoPrometheus := echoPrometheus.NewPrometheus("statuspage_exporter", nil)
  echoPrometheus.Use(srv)

  srv.GET("/probe", prober.Handler(log))
  srv.GET("/healthz", handleHealthz)

  httpPort := config.HTTPPort()
  httpAddr := fmt.Sprintf(":%d", httpPort)

  // Start your http server for prometheus.
  go func() {
    if err := srv.Start(httpAddr); err != nil {
      log.Panic("Unable to start a http server.", zap.Error(err))
    }
  }()

  log.Info("Http server listening on", zap.Int("port", httpPort))

  <-ctx.Done()

  shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
  defer cancel()

  if err := srv.Shutdown(shutdownCtx); err != nil { //nolint:contextcheck
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

  go startHTTP(ctx, wg, log)

  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

  <-quit
  log.Info("Received shutdown signal. Waiting for workers to terminate...")
  cancel()

  wg.Wait()
}
