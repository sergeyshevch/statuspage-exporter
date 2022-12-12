package core

import (
	"context"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/sergeyshevch/statuspage-exporter/pkg/config"
)

// StartFetchingLoop starts a loop that fetches status pages.
func StartFetchingLoop(ctx context.Context, wg *sync.WaitGroup,
	log *zap.Logger, client *resty.Client, targetURLs []string,
	fetcher func(*zap.Logger, *resty.Client, []string),
) {
	wg.Add(1)
	defer wg.Done()

	fetchDelay := config.FetchDelay()

	for {
		select {
		default:
			fetcher(log, client, targetURLs)

			time.Sleep(fetchDelay)
		case <-ctx.Done():
			log.Info("Stopping fetching loop")

			return
		}
	}
}
