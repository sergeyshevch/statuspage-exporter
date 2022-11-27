package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	defaultClientTimeout = 2 * time.Second
	defaultFetchDelay    = 5 * time.Second
	defaultHTTPPort      = 8080
)

var configMutex = &sync.Mutex{}

// InitConfig initializes a config and configure viper to receive config from file and environment.
func InitConfig() (*zap.Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to create logger", zap.Error(err))
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return log, err
	}

	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	viper.SetConfigName(".statuspage-exporter")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}

	return log, nil
}

// HTTPPort returns a port for http server.
func HTTPPort() int {
	viper.SetDefault("http_port", defaultHTTPPort)

	return viper.GetInt("http_port")
}

// FetchDelay returns a delay between statuspages fetches.
func FetchDelay() time.Duration {
	configMutex.Lock()
	viper.SetDefault("fetch_delay", defaultFetchDelay)
	value := viper.GetDuration("fetch_delay")
	configMutex.Unlock()

	return value
}

// ClientTimeout returns a timeout for http client.
func ClientTimeout() time.Duration {
	configMutex.Lock()
	viper.SetDefault("client_timeout", defaultClientTimeout)
	value := viper.GetDuration("client_timeout")
	configMutex.Unlock()

	return value
}

// StatusPages returns a list of status pages to monitor.
func StatusPages() []string {
	configMutex.Lock()
	value := viper.GetStringSlice("statuspages")
	configMutex.Unlock()

	return value
}
