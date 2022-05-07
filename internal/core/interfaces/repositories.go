package interfaces

import (
	"context"
	"service-area-service/internal/core/domain"
)

type ServiceAreaRepository interface {
	GetAll(ctx context.Context) ([]domain.ServiceArea, error)
	Get(ctx context.Context, id int) (domain.ServiceArea, error)
	Save(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error)
	Update(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error)
}
