package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

type Option func(logger *Logger)

type Logger struct {
	resolvers map[string]func(*gin.Context) string
	notLogged map[string]struct{}
}

// New - Gets a new logger
func New(opts ...Option) Logger {
	logger := Logger{
		resolvers: make(map[string]func(*gin.Context) string),
		notLogged: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(&logger)
	}
	return logger
}

// SetResolver - Sets a given field and the resolver to resolve the field
func (r Logger) SetResolver(name string, resolver func(*gin.Context) string) {
	r.resolvers[name] = resolver
}

// SetNotLogged - Sets the paths that should not be logged
func (r Logger) SetNotLogged(notLogged ...string) {
	for _, path := range notLogged {
		r.notLogged[path] = struct{}{}
	}
}

// Middleware - Gets the Gin Middleware logger
func (r Logger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := getDurationInMilliseconds(start)

		if _, ok := r.notLogged[c.Request.RequestURI]; !ok {
			// Define standard records
			ctx := log.With().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", c.Request.RequestURI).
				Str("client_ip", getClientIP(c)).
				Float64("duration", duration).
				Str("referrer", c.Request.Referer()).
				Str("user_agent", c.Request.UserAgent()).
				Str("trace_id", getTraceID(c))

			// Resolve custom entities
			if len(r.resolvers) > 0 {
				for name, resolve := range r.resolvers {
					ctx = ctx.Str(name, resolve(c))
				}
			}

			logger := ctx.Logger()

			if c.Writer.Status() >= 500 {
				logger.Error().Msg(c.Errors.String())
			} else {
				logger.Info().Send()
			}
		}
	}
}

func getDurationInMilliseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}

func getTraceID(c *gin.Context) string {
	span := trace.SpanFromContext(c.Request.Context())
	var spanID, traceID string
	if span.SpanContext().HasSpanID() && span.SpanContext().SpanID().IsValid() {
		spanID = span.SpanContext().SpanID().String()
	}
	if span.SpanContext().HasTraceID() && span.SpanContext().TraceID().IsValid() {
		traceID = span.SpanContext().TraceID().String()
	}

	return fmt.Sprintf("1-%s-%s", traceID, spanID)
}

func getClientIP(c *gin.Context) string {
	// first check the X-Forwarded-For header
	requester := c.Request.Header.Get("X-Forwarded-For")
	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}
