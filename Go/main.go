package main

import (
	_ "embed"
	ejson "encoding/json"
	"example.com/kafka-serializer-publisher/api"
	"example.com/kafka-serializer-publisher/config"
	"example.com/kafka-serializer-publisher/config/openapi"
	"example.com/kafka-serializer-publisher/config/router"
	"example.com/kafka-serializer-publisher/config/telemetry"
	"example.com/kafka-serializer-publisher/kafka"
	"example.com/kafka-serializer-publisher/marshaller"
	"example.com/kafka-serializer-publisher/marshaller/avro"
	"example.com/kafka-serializer-publisher/marshaller/json"
	"example.com/kafka-serializer-publisher/model"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

//go:embed metadata
var src string

// @title        Kafka Publisher API
// @version      1.0
// @description  API that publishes messages to a kafka message broker

// @contact.name   Sean Doyle
// @contact.email  doyle316@umn.edu

// @host      localhost:8080
// @BasePath  /
// @schemes   http https
func main() {
	// Set Configuration
	cfg := config.Load(false, "./defaults.yaml", "./.env")
	arr, _ := ejson.Marshal(cfg)
	fmt.Printf("Config %s", string(arr))

	// Start Event Producer
	marsh := marshaller.New(json.New(), avro.New(cfg.SchemaRegistry))
	sender := kafka.New(cfg.Kafka, marsh, kafka.WithTracing(cfg.OTEL.Enabled))

	// Setup Routes & Metrics, Start Server
	err := startServer(cfg, sender)
	if err != nil {
		log.Fatal().Msgf("Error setting routes: %s", err)
		return
	}
}

// startServer - Sets up API
func startServer(cfg *config.Config, sender model.Publisher) error {
	options := []router.Option{
		router.WithDebug(cfg.App.Debug),
		router.WithHealth("/health"),
		router.WithMetrics("/metrics"),
		router.WithPathExclusions(openapi.GetIgnoredRoutes()...),
		router.WithTrustedProxies(cfg.Proxies),
	}
	// Start Router & Tracing
	if cfg.OTEL.Enabled {
		options = append(options, router.WithTraceProvider(cfg.App.Name, telemetry.New(cfg.App.Name, cfg.OTEL)))
	}
	rtr := router.New(options...)

	// Configure Routes
	openapi.Setup(cfg, rtr.Engine)
	api.CreateRoutes(rtr.Engine, sender)

	// Start Server
	hostname, _ := os.Hostname()
	log.Debug().Msgf("Starting server on the port http://%s:%d", hostname, cfg.App.Port)
	fmt.Println(src)
	return rtr.Engine.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
