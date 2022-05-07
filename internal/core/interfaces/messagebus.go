package interfaces

import (
	"context"
	"service-area-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error
	UpdateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error
}
