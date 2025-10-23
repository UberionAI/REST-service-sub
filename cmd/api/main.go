package main

import (
	_ "REST-service-sub/docs"
	"REST-service-sub/internal/config"
	"REST-service-sub/internal/db"
	"REST-service-sub/internal/handler"
	"REST-service-sub/internal/middleware"
	"REST-service-sub/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title REST API Subscription service
// @version 1.0
// @description REST API for userdata aggregation of the subscriptions
// @host localhost:8000
// @BasePath /

func main() {
	_ = godotenv.Load("../../.env")
	cfg := config.LoadConfig()
	fmt.Println("Config loaded successfully...")

	//Checking DSN (check for getting .env from docker)
	//	fmt.Println("DSN:", cfg.DSN())

	gdb, err := db.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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
		log.Fatalf("Failed to start server: %v", err)
	}
}
