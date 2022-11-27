package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

import (
	"github.com/spf13/viper"
)

var configMutex = &sync.Mutex{}

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

func HTTPPort() int {
	return 8000
}

func FetchDelay() time.Duration {
	return 5 * time.Second
}

func ClientTimeout() time.Duration {
	return 2 * time.Second
}

func StatusPages() []string {
	configMutex.Lock()
	value := viper.GetStringSlice("statuspages")
	configMutex.Unlock()
	return value
}
