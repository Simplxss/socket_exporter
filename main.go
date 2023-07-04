package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-yaml/yaml"

	"github.com/YuFanXing/socket_exporter/client"
	"github.com/YuFanXing/socket_exporter/config"
)

var clients []client.Client

func main() {
	var config config.Config

	data, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("cannot open config file: %v", err)
	}
	err = yaml.NewDecoder(data).Decode(&config)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	reg := prometheus.NewRegistry()

	for _, endpoint := range config.Endpoint {
		clients = append(clients, *client.NewClient(endpoint, reg))
	}

	http.Handle(config.Listening.Path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	log.Fatal(http.ListenAndServe(config.Listening.Address, nil))
}
