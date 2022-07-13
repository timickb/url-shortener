package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/timickb/url-shortener/internal/app/urlapi"
	"gopkg.in/yaml.v3"
)

var (
	configPath string
	storeImpl  string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/server.yaml", "Path to config file")
	flag.StringVar(&storeImpl, "store", "local", "Data storage way")
}

func main() {
	flag.Parse()

	config := urlapi.DefaultConfig()

	configContent, ioErr := ioutil.ReadFile(configPath)

	if ioErr != nil {
		log.Fatal(ioErr)
	}

	yamlErr := yaml.Unmarshal(configContent, &config)

	if yamlErr != nil {
		log.Fatal(yamlErr)
	}

	server := urlapi.NewServer(config)

	err := server.Start(storeImpl)

	if err != nil {
		log.Fatal(err)
	}
}
