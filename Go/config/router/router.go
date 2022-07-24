package router

import (
	"example.com/kafka-serializer-publisher/config/logging"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Option func(*Router)

type Router struct {
	Engine      *gin.Engine
	logger      *logging.Logger
	monExcludes map[string]string
}

func New(opts ...Option) Router {
	// Tell gin to use Zerolog as its logging framework
	gin.DefaultWriter = log.Logger

	router := Router{
		Engine:      gin.New(),
		monExcludes: make(map[string]string),
	}
	router.Engine.RedirectTrailingSlash = false

	// Setup Logging Details
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	WithDebug(false)(&router)

	// Loop through each option
	for _, opt := range opts {
		opt(&router)
	}

	// Tell logger to not log requests to specified paths
	var exclusions []string
	for k := range router.monExcludes {
		exclusions = append(exclusions, k)
	}

	logger := logging.New(logging.WithExclusions(exclusions...))
	router.logger = &logger
	router.Engine.Use(router.logger.Middleware())

	return router
}

// IsExcludedOTEL Checks if path is excluded
func (a Router) IsExcludedOTEL(path string) bool {
	_, ok := a.monExcludes[path]
	return ok
}
