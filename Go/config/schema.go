package config

import (
	"example.com/kafka-serializer-publisher/config/telemetry"
	"example.com/kafka-serializer-publisher/kafka"
	"example.com/kafka-serializer-publisher/marshaller/avro"
)

type Config struct {
	App            App                 `yaml:"app"`
	Proxies        []string            `yaml:"proxies"`
	Kafka          kafka.Config        `yaml:"kafka"`
	OTEL           telemetry.Config    `yaml:"otel" mapstructure:"otel"`
	SchemaRegistry avro.SchemaRegistry `yaml:"schema_registry" mapstructure:"schema_registry"`
}

type App struct {
	Name     string `yaml:"name"`
	Debug    bool   `yaml:"debug"`
	Port     int    `yaml:"port"`
	BasePath string `yaml:"base_path" mapstructure:"base_path"`
}
