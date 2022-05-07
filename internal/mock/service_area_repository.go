package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"service-area-service/internal/core/domain"
)

type ServiceAreaRepository struct {
	mock.Mock
}

func (m *ServiceAreaRepository) GetAll(ctx context.Context) ([]domain.ServiceArea, error) {
	args := m.Called()
	return args.Get(0).([]domain.ServiceArea), args.Error(1)
}

func (m *ServiceAreaRepository) Get(ctx context.Context, id int) (domain.ServiceArea, error) {
	args := m.Called(id)
	return args.Get(0).(domain.ServiceArea), args.Error(1)
}

func (m *ServiceAreaRepository) Save(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error) {
	args := m.Called(serviceArea)
	return args.Get(0).(domain.ServiceArea), args.Error(1)
}

func (m *ServiceAreaRepository) Update(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error) {
	args := m.Called(serviceArea)
	return args.Get(0).(domain.ServiceArea), args.Error(1)
}
