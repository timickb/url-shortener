package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/timickb/url-shortener/internal/app/algorithm"
	"github.com/timickb/url-shortener/internal/app/server"
	"github.com/timickb/url-shortener/internal/app/store"
	"gopkg.in/yaml.v3"

	_ "github.com/lib/pq"
)

var (
	configSource string
	configPath   string
)

const (
	defaultConfigPath = "config.yml"
)

func init() {
	flag.StringVar(&configPath, "config-path", "config.yml", "path to config in filesystem")
	flag.StringVar(&configSource, "config-source", "env", "config params source priority: env or file")
}

func main() {
	logger := logrus.New()
	flag.Parse()

	if err := mainNoExit(logger); err != nil {
		log.Fatalf("fatal err: %s", err.Error())
	}
}

func mainNoExit(logger *logrus.Logger) error {
	config := server.DefaultConfig()

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

	fmt.Println(config.Origins)

	db, err := createPostgresConn(config)

	if err != nil {
		logger.Warnf("Couldnt create database connection: %s", err.Error())
		logger.Warnf("API will use only in-memory storage!")
	} else {
		logger.Info("Database connection established")
	}

	shr := algorithm.DefaultShortener{HashSize: config.ShorteningSize}

	st := store.New(
		store.WithDB(db),
		store.WithLogger(logger),
		store.WithShortener(shr),
	)

	if err := st.Open(); err != nil {
		return err
	}

	store := st
	logger.Info("Store configured")

	// Creating server
	srv := server.New(
		server.WithConfig(config),
		server.WithLogger(logger),
		server.WithStore(store),
	)

	if err := srv.Start(); err != nil {
		return err
	}

	return nil
}

func parseConfigFromFile(config *server.Config) {
	if configPath == "" {
		configPath = strings.Clone(defaultConfigPath)
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
	if os.Getenv("SHORTENING_SIZE") != "" {
		config.ShorteningSize, _ = strconv.Atoi(os.Getenv("SHORTENING_SIZE"))
	}
}

func createPostgresConn(config *server.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=\"%s\" dbname=%s sslmode=disable",
		config.Database.DbHost,
		config.Database.DbPort,
		config.Database.DbUser,
		config.Database.DbPassword,
		config.Database.DbName)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}
