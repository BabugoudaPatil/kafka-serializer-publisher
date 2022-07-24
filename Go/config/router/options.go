package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	sdktrace "go.opentelemetry.io/otel/trace"
)

// WithTraceProvider - Enables monitoring with Gin and OTEL
func WithTraceProvider(serviceName string, tracing sdktrace.TracerProvider) Option {
	return func(r *Router) {
		middleware := otelgin.Middleware(serviceName,
			otelgin.WithTracerProvider(tracing),
			otelgin.WithPropagators(xray.Propagator{}),
		)

		r.Engine.Use(func(c *gin.Context) {
			if !r.IsExcludedOTEL(c.FullPath()) {
				middleware(c)
			}
		})
	}
}

func WithDebug(enabled bool) Option {
	return func(r *Router) {
		if enabled {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			gin.SetMode(gin.DebugMode)
			pprof.Register(r.Engine)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			gin.SetMode(gin.ReleaseMode)
		}
	}
}

// WithMetrics - Enables monitoring with Gin and Promethous
func WithMetrics(path string) Option {
	return func(r *Router) {
		// Setup Prometheus Metrics
		g := ginmetrics.GetMonitor()
		g.SetMetricPath(path)
		// Set slow time, default 5s
		g.SetSlowTime(10)
		// Set request duration, default {0.1, 0.3, 1.2, 5, 10}
		// used to p95, p99
		g.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		// set Promethous metrics for gin
		g.Use(r.Engine)
		// Disable logging of metrics
		WithPathExclusions(path)(r)
	}
}

// WithHealth - Enables health endpoints
func WithHealth(path string) Option {
	return func(r *Router) {
		r.Engine.GET(path, func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "active",
			})
		})
		// Disable logging of health checks
		WithPathExclusions(path)(r)
	}
}

// WithPathExclusions - Excludes multiple paths from sending to OTEL
func WithPathExclusions(names ...string) Option {
	return func(r *Router) {
		for _, name := range names {
			r.monExcludes[name] = ""
		}
	}
}

// WithTrustedProxies - Excludes multiple paths from sending to OTEL
func WithTrustedProxies(names []string) Option {
	return func(r *Router) {
		if err := r.Engine.SetTrustedProxies(names); err != nil {
			log.Fatal().Stack().Err(err).Msgf("Failed to set trusted proxies")
		}
	}
}
