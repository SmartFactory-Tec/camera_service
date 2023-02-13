package main

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path"
)

type (
	Config struct {
		Port int

		Db DbConfig
	}
)

func loadConfig(logger *zap.SugaredLogger) Config {
	configLoader := viper.New()

	if v, ok := os.LookupEnv("CAMERA_SERVICE_CONFIG"); ok {
		configLoader.AddConfigPath(v)
		if err := os.MkdirAll(v, 0777); err != nil {
			logger.Fatal(err)
		}
	} else if v, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		folder := path.Join(v, "camera_service")
		if err := os.MkdirAll(folder, 0777); err != nil {
			logger.Fatal(err)
		}
		configLoader.AddConfigPath(folder)
	} else if v, ok := os.LookupEnv("HOME"); ok {
		folder := path.Join(v, ".config/camera_service")
		if err := os.MkdirAll(folder, 0777); err != nil {
			logger.Fatal(err)
		}
		configLoader.AddConfigPath(folder)
	} else {
		logger.Fatal("could not resolve config path")
	}

	configLoader.SetConfigName("config")
	configLoader.SetConfigType("toml")

	// main config
	configLoader.SetDefault("port", 3000)

	// db config
	configLoader.SetDefault("db.hostname", "localhost")
	configLoader.SetDefault("db.port", 5432)
	configLoader.SetDefault("db.database", "")
	configLoader.SetDefault("db.user", "")
	configLoader.SetDefault("db.password", "")

	err := configLoader.ReadInConfig()

	if err != nil {
		var notFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &notFoundError) {
			logger.Infow("file not found, creating default config file", "file", configLoader.ConfigFileUsed())
			if err := configLoader.SafeWriteConfig(); err != nil {
				logger.Fatal(fmt.Errorf("could not create default config file: %w", err))
			}
		} else {
			logger.Fatal(fmt.Errorf("could not read config file: %w", err))
		}
	} else {
		logger.Infow("loaded service config from config file", "file", configLoader.ConfigFileUsed())
	}

	var config Config

	if err := configLoader.Unmarshal(&config); err != nil {
		logger.Fatal(fmt.Errorf("error unmarshaling config file: %w", err))
		return Config{}
	}

	return config
}
