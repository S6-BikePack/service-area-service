package main

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"service-area-service/config"
	"service-area-service/internal/core/services"
	"service-area-service/internal/handlers"
	"service-area-service/internal/repositories"
	"service-area-service/pkg/logging"
	"service-area-service/pkg/rabbitmq"
	"service-area-service/pkg/tracing"

	"github.com/gin-gonic/gin"
)

const defaultConfig = "./config/local.config"

func main() {
	cfgPath := GetEnvOrDefault("config", defaultConfig)
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(err)
	}

	logger, err := logging.NewSugaredOtelZap(cfg)
	defer logger.Close()

	if err != nil {
		panic(err)
	}

	tracer, err := tracing.NewOpenTracing(cfg.Server.Service, cfg.Tracing.Host, cfg.Tracing.Port)

	//--------------------------------------------------------------------------------------
	// Setup Database
	//--------------------------------------------------------------------------------------

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	serviceAreaRepository, err := repositories.NewServiceAreaRepository(db)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	//--------------------------------------------------------------------------------------
	// Setup RabbitMQ
	//--------------------------------------------------------------------------------------

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	rmqPublisher := services.NewRabbitMQPublisher(rmqServer, tracer, cfg)

	//--------------------------------------------------------------------------------------
	// Setup Services
	//--------------------------------------------------------------------------------------

	serviceAreaService := services.NewServiceAreaService(cfg, rmqPublisher, serviceAreaRepository)

	//--------------------------------------------------------------------------------------
	// Setup HTTP server
	//--------------------------------------------------------------------------------------

	router := gin.New()

	deliveryHandler := handlers.NewHTTPHandler(serviceAreaService, router, logger, cfg)
	deliveryHandler.SetupEndpoints()
	deliveryHandler.SetupSwagger()

	logger.Fatal(context.Background(), router.Run(cfg.Server.Port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	fmt.Println(returnValue)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
