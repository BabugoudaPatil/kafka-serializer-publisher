package telemetry

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	Enabled    bool          `yaml:"enabled"`
	Host       string        `yaml:"host"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxMsgSize int           `yaml:"max_msg_size" mapstructure:"max_msg_size"`
}

func (c Config) getMaxMsgSize() int {
	if c.MaxMsgSize == 0 {
		return 2000000
	}
	return c.MaxMsgSize
}

func (c Config) getTimeout() time.Duration {
	if c.Timeout == time.Second*0 {
		return time.Second * 10
	}
	return c.Timeout
}

type TraceProvider struct {
	*sdktrace.TracerProvider
}

// New - Initialize tracing
func New(serviceName string, cfg Config) TraceProvider {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.getTimeout())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("failed to create resource")
	}

	// GRPC Connect to OTEL Collector
	conn, err := grpc.DialContext(ctx, cfg.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(cfg.getMaxMsgSize())))
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to create gRPC connection to collector")
	}

	// Setup X-Ray ID
	idg := xray.NewIDGenerator()

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("Failed to create OTEL trace exporter")
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(traceExporter)),
		sdktrace.WithIDGenerator(idg),
	)
	// Set global trace provider
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(xray.Propagator{})

	return TraceProvider{tracerProvider}
}

func (t *TraceProvider) Disconnect() error {
	if err := t.Shutdown(context.Background()); err != nil {
		log.Debug().Msgf("Error shutting down tracer provider: %v", err)
		return err
	}
	return nil
}
