package kafka

import (
	"context"
	"errors"
	"example.com/kafka-serializer-publisher/marshaller"
	"example.com/kafka-serializer-publisher/model"
	"fmt"
	"github.com/Shopify/sarama"
	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog/log"
	"time"
)

type Option func(*producer)

type producer struct {
	cfg           Config
	converter     marshaller.Marshaller
	sender        sarama.SyncProducer
	enableTracing bool
}

func New(cfg Config, marsh marshaller.Marshaller, options ...Option) model.Publisher {
	log.Debug().Msg("Starting Kafka Producer")
	p := producer{cfg: cfg, converter: marsh}
	for _, opt := range options {
		opt(&p)
	}
	var err error
	p.sender, err = sarama.NewSyncProducer(cfg.Hosts, cfg.SaramaConfig(p.enableTracing))
	if err != nil {
		log.Fatal().Msgf("failed to create producer: %s", err.Error())
	}
	return &p
}

func (p *producer) Send(ctx context.Context, contentType model.EventContentType, publishable model.Publishable) error {
	_ = ctx
	msg, err := p.buildMessage(publishable, contentType)
	if err != nil {
		return err
	}

	prometheusClient := prometheusmetrics.NewPrometheusProvider(
		metrics.DefaultRegistry, "", "producer", prometheus.DefaultRegisterer, 1*time.Second)
	go prometheusClient.UpdatePrometheusMetrics()

	_, _, err = p.sender.SendMessage(&msg)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to send event : %s", err.Error())
		return errors.New(fmt.Sprintf("Failed to send event : %v", err.Error()))
	}
	return nil
}

func (p *producer) buildMessage(publishable model.Publishable, contentType model.EventContentType) (sarama.ProducerMessage, error) {
	data, err := p.converter.Marshal(publishable, contentType)
	if err != nil {
		log.Error().Err(err)
		return sarama.ProducerMessage{}, err
	}
	var headers []sarama.RecordHeader
	for key, value := range publishable.Headers() {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}
	headers = append(headers, sarama.RecordHeader{
		Key:   []byte("content-type"),
		Value: []byte(contentType),
	})

	msg := sarama.ProducerMessage{
		Topic:     publishable.Topic(),
		Value:     sarama.ByteEncoder(data),
		Headers:   headers,
		Timestamp: time.Now(),
	}
	if publishable.ID() != "" {
		msg.Key = sarama.StringEncoder(publishable.ID())
	}
	return msg, nil
}

// Disconnect - Disconnects Producer connection
func (p *producer) Disconnect() {
	log.Debug().Msg("Closing Kafka Producer")
	if err := p.sender.Close(); err != nil {
		log.Fatal().Msgf("failed to close writer: %s", err)
	}
}
