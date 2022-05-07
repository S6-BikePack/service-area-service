package interfaces

import (
	"context"
	"service-area-service/internal/core/domain"
)

type ServiceAreaService interface {
	GetAll(ctx context.Context) ([]domain.ServiceArea, error)
	Get(ctx context.Context, id int) (domain.ServiceArea, error)
	Create(ctx context.Context, id int, identifier, name string, area domain.Area) (domain.ServiceArea, error)
	Update(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error)
}
