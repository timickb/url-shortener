package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/server"
	"gopkg.in/yaml.v3"
)

var (
	configSource string
	configPath   string
	storeImpl    string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.yml", "path to config in filesystem")
	flag.StringVar(&configSource, "config-source", "env", "config params source priority: env or file")
	flag.StringVar(&storeImpl, "store", "local", "Data storage way: local, db or improved")
}

func main() {
	logger := logrus.New()

	if err := mainNoExit(logger); err != nil {
		log.Fatalf("fatal err: %s", err.Error())
	}
}

func mainNoExit(logger *logrus.Logger) error {
	flag.Parse()

	config := server.DefaultConfig()

	config.StoreImpl = storeImpl

	// Parsing configuration
	if configSource == "file" {
		logger.Info("Reading configuration from file")
		parseConfigFromFile(config)
	} else if configSource == "env" {
		logger.Info("Reading configuration from environment")
		parseConfigFromEnvironment(config)
	} else {
		return fmt.Errorf("incorrect config source %s. Use 'file' or 'env'", configSource)
	}

	// Creating server
	srv, err := server.New(config, logger)

	if err != nil {
		return err
	}

	if err := srv.Start(); err != nil {
		return err
	}

	return nil
}

func parseConfigFromFile(config *server.Config) {
	if configPath == "" {
		configPath = "config.yml"
	}

	configContent, ioErr := ioutil.ReadFile(configPath)

	if ioErr != nil {
		log.Fatal(ioErr)
	}

	if err := yaml.Unmarshal(configContent, &config); err != nil {
		log.Fatal(err)
	}
}

func parseConfigFromEnvironment(config *server.Config) {
	config.Database.DbHost = os.Getenv("DB_HOST")
	config.Database.DbUser = os.Getenv("DB_USER")
	config.Database.DbName = os.Getenv("DB_NAME")

	config.Database.DbPassword = os.Getenv("DB_PASSWORD")
	config.Database.DbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))

	if os.Getenv("SERVER_PORT") != "" {
		config.ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	}
	if os.Getenv("MAX_URL_LENGTH") != "" {
		config.MaxUrlLength, _ = strconv.Atoi(os.Getenv("MAX_URL_LENGTH"))
	}

	if os.Getenv("STORE_WAY") != "" {
		config.StoreImpl = os.Getenv("STORE_WAY")
	}
}
