package handlers

import (
	"fmt"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"service-area-service/config"
	"service-area-service/internal/core/interfaces"
	"service-area-service/pkg/authorization"
	"service-area-service/pkg/dto"
	"service-area-service/pkg/logging"
	"strconv"

	ginSwagger "github.com/swaggo/gin-swagger"
	"service-area-service/docs"
	_ "service-area-service/docs"
)
import "github.com/gin-gonic/gin"

type HTTPHandler struct {
	serviceAreaService interfaces.ServiceAreaService
	router             *gin.Engine
	logger             logging.Logger
	config             *config.Config
}

func NewHTTPHandler(serviceAreaService interfaces.ServiceAreaService, router *gin.Engine, logger logging.Logger, config *config.Config) *HTTPHandler {
	return &HTTPHandler{
		serviceAreaService: serviceAreaService,
		router:             router,
		logger:             logger,
		config:             config,
	}
}

func (handler *HTTPHandler) SetupEndpoints() {
	api := handler.router.Group("/api")
	api.GET("/service-areas", handler.GetAll)
	api.GET("/service-areas/:id", handler.Get)
	api.POST("/service-areas", handler.Create)
}

func (handler *HTTPHandler) SetupSwagger() {
	docs.SwaggerInfo.Title = handler.config.Server.Service + " API"
	docs.SwaggerInfo.Description = handler.config.Server.Description

	handler.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// GetAll godoc
// @Summary  get all service-areas
// @Schemes
// @Description  gets all service-areas in the system
// @Accept       json
// @Produce      json
// @Success      200  {object}  []dto.ServiceAreaListResponse
// @Router       /api/service-areas [get]
func (handler *HTTPHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	serviceAreas, err := handler.serviceAreaService.GetAll(ctx)

	if err != nil {
		handler.logger.Error(ctx, err.Error(), "error", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, dto.CreateServiceAreaListResponse(serviceAreas))
}

// Get godoc
// @Summary  get service-area
// @Schemes
// @Param        id     path  string           true  "Service-area id"
// @Description  gets a service-area from the system by its ID
// @Produce      json
// @Success      200  {object}  dto.ServiceAreaResponse
// @Router       /api/service-areas/{id} [get]
func (handler *HTTPHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	serviceArea, err := handler.serviceAreaService.Get(ctx, id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, dto.CreateServiceAreaResponse(serviceArea))
}

// Create godoc
// @Summary  create service-area
// @Schemes
// @Description  creates a new service-area
// @Accept       json
// @Param        service-area  body  dto.BodyCreateServiceArea  true  "Add service-area"
// @Produce      json
// @Success      200  {object}  dto.ServiceAreaResponse
// @Router       /api/service-areas [post]
func (handler *HTTPHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	body := dto.BodyCreateServiceArea{}
	err := c.BindJSON(&body)

	if body.ID == 0 || body.Name == "" || body.Identifier == "" || err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	auth := authorization.NewRest(c)

	if auth.AuthorizeAdmin() || true {

		serviceArea, err := handler.serviceAreaService.Create(ctx, body.ID, body.Identifier, body.Name, body.Area)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Could not create service area with the folling data: %v", body)})
			return
		}

		c.JSON(http.StatusOK, dto.CreateServiceAreaResponse(serviceArea))
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}
