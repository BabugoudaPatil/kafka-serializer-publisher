package openapi

import (
	"example.com/kafka-serializer-publisher/config"
	"example.com/kafka-serializer-publisher/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(cfg *config.Config, router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configure API Docs
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{}
	if cfg.App.BasePath != "" {
		docs.SwaggerInfo.BasePath = cfg.App.BasePath
	}
}

func GetIgnoredRoutes() []string {
	return []string{"/swagger/swagger-ui.css", "/swagger/swagger-ui-standalone-preset.js", "/swagger/swagger-ui-bundle.js", "/swagger/doc.json", "/swagger/favicon-32x32.png", "/favicon.ico"}
}
