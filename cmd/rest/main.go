package main

import (
	"context"
	"fmt"
	"service-area-service/pkg/azure"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"os"
	"service-area-service/config"
	"service-area-service/internal/core/services"
	"service-area-service/internal/handlers"
	"service-area-service/internal/repositories"
	"service-area-service/pkg/logging"
	"service-area-service/pkg/tracing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const defaultConfig = "./config/local.config"

func main() {
	cfgPath := GetEnvOrDefault("config", defaultConfig)
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		fmt.Printf("Failed to load config: %v", err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Logging and Tracing
	//--------------------------------------------------------------------------------------

	logger, err := logging.NewSimpleLogger(cfg)

	if err != nil {
		panic(err)
	}

	tracer, err := tracing.NewOpenTracing(cfg.Server.Service, cfg.Tracing.Host, cfg.Tracing.Port)

	if err != nil {
		logger.Warning(context.Background(), "Failed to setup tracing: %v", err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Database
	//--------------------------------------------------------------------------------------

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database, cfg.Database.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	if tracer != nil {
		if err = db.Use(otelgorm.NewPlugin(otelgorm.WithTracerProvider(tracer))); err != nil {
			logger.Fatal(context.Background(), err)
		}
	}

	serviceAreaRepository, err := repositories.NewServiceAreaRepository(db)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Azure Service Bus
	//--------------------------------------------------------------------------------------

	azServiceBus, err := azure.NewAzureServiceBus(cfg)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	azPublisher := services.NewAzurePublisher(azServiceBus, cfg)

	//--------------------------------------------------------------------------------------
	// Setup Services
	//--------------------------------------------------------------------------------------

	serviceAreaService := services.NewServiceAreaService(cfg, azPublisher, serviceAreaRepository)

	//--------------------------------------------------------------------------------------
	// Setup HTTP server
	//--------------------------------------------------------------------------------------

	router := gin.New()

	if tracer != nil {
		router.Use(otelgin.Middleware(cfg.Server.Service, otelgin.WithTracerProvider(tracer)))
	}

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
