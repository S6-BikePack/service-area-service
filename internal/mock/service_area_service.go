package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"service-area-service/internal/core/domain"
)

type ServiceAreaService struct {
	mock.Mock
}

func (m *ServiceAreaService) GetAll(ctx context.Context) ([]domain.ServiceArea, error) {
	args := m.Called()
	return args.Get(0).([]domain.ServiceArea), args.Error(1)
}

func (m *ServiceAreaService) Get(ctx context.Context, id int) (domain.ServiceArea, error) {
	args := m.Called(id)
	return args.Get(0).(domain.ServiceArea), args.Error(1)
}

func (m *ServiceAreaService) Create(ctx context.Context, id int, identifier, name string, area domain.Area) (domain.ServiceArea, error) {
	args := m.Called(id, identifier, name, area)
	return args.Get(0).(domain.ServiceArea), args.Error(1)
}

func (m *ServiceAreaService) Update(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error) {
	args := m.Called(serviceArea)
	return args.Get(0).(domain.ServiceArea), args.Error(1)
}
