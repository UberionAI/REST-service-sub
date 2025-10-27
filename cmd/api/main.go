package main

import (
	_ "REST-service-sub/docs"
	"REST-service-sub/internal/config"
	"REST-service-sub/internal/db"
	"REST-service-sub/internal/handler"
	"REST-service-sub/internal/logger"
	"REST-service-sub/internal/middleware"
	"REST-service-sub/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strings"
)

// @title REST API Subscription service
// @version 1.0
// @description REST API for userdata aggregation of the subscriptions
// @host localhost:8000
// @BasePath /

func main() {
	_ = godotenv.Load("../../.env")
	cfg := config.LoadConfig()

	logger.Init(cfg.LogLevel)
	log.Info().Str("log_level", cfg.LogLevel).Msg("Logger is initialized")

	fmt.Println("Config loaded successfully...")

	if strings.ToLower(cfg.LogLevel) == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gdb, err := db.NewGormDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	fmt.Println("Database connected successfully!")

	subService := service.NewSubscriptionService(gdb)
	subHandler := handler.NewSubscriptionHandler(subService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// Swagger init
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, "OK")
	})

	subHandler.RegisterRoutes(r)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("Server is listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
