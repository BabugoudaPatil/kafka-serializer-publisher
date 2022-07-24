package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/Shopify/sarama"
	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"software.sslmate.com/src/go-pkcs12"
	"time"
)

type Config struct {
	Hosts            []string `yaml:"hosts"`
	Source           string   `yaml:"source"`
	CACertPath       string   `yaml:"ca_cert_path" mapstructure:"ca_cert_path"`
	KeystorePath     string   `yaml:"keystore_path" mapstructure:"keystore_path"`
	KeystorePassword string   `yaml:"keystore_password" mapstructure:"keystore_password"`
	AllowCreateNew   bool     `yaml:"allow_create_new" mapstructure:"allow_create_new"`
}

// SaramaConfig - Gets basic Sarama configuration for Kafka
func (k Config) SaramaConfig(withMetrics bool) *sarama.Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Version = sarama.V2_3_1_0
	saramaConfig.Metadata.AllowAutoTopicCreation = k.AllowCreateNew
	if withMetrics {
		metricGenerator(saramaConfig)
	}
	tlsConfig, err := k.tlsConfig()
	if k.IsTLSConfigured() {
		if err != nil {
			log.Info().Stack().Err(err).Msgf("Kafka TLS authentication is configured but unavailable (CA:%s, Key:%s)", k.CACertPath, k.KeystorePath)
		} else {
			log.Info().Msgf("Enabling secure kafka connection via TLS (CA:%s, Key:%s)", k.CACertPath, k.KeystorePath)
			saramaConfig.Net.TLS.Enable = true
			saramaConfig.Net.TLS.Config = tlsConfig
		}
	}
	return saramaConfig
}

func (k Config) IsTLSConfigured() bool {
	if k.CACertPath == "" || k.KeystorePath == "" || k.KeystorePassword == "" {
		return false
	}
	return true
}

func (k Config) tlsConfig() (*tls.Config, error) {
	// Keystore
	cert, err := k.keystoreCert()
	if err != nil {
		return nil, err
	}

	//Truststore
	if k.CACertPath == "" {
		return nil, errors.New("CA cert path not set")
	}

	caCert, err := ioutil.ReadFile(k.CACertPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		RootCAs:      caCertPool,
	}
	return config, nil
}

func (k Config) keystoreCert() (*tls.Certificate, error) {
	if k.KeystorePath == "" {
		return nil, errors.New("keystore path not set")
	} else if k.KeystorePassword == "" {
		return nil, errors.New("keystore password not set")
	}

	data, err := ioutil.ReadFile(k.KeystorePath)
	if err != nil {
		return nil, err
	}
	pk, crt, _, err := pkcs12.DecodeChain(data, k.KeystorePassword)
	if err != nil {
		return nil, err
	}
	tlsCrt := tls.Certificate{
		Certificate: [][]byte{crt.Raw},
		Leaf:        crt,
		PrivateKey:  pk,
	}
	return &tlsCrt, nil
}

func metricGenerator(saramaConfig *sarama.Config) {
	appMetricRegistry := saramaConfig.MetricRegistry
	saramaConfig.MetricRegistry = metrics.NewPrefixedChildRegistry(appMetricRegistry, "kafka.")
	prometheusClient := prometheusmetrics.NewPrometheusProvider(
		appMetricRegistry, "", "", prometheus.DefaultRegisterer, 1*time.Second)
	go prometheusClient.UpdatePrometheusMetrics()
	http.Handle("/metrics", promhttp.Handler())
}
