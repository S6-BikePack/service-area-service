package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"service-area-service/internal/core/domain"
)

type RabbitMQPublisher struct {
	mock.Mock
}

func (m *RabbitMQPublisher) CreateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error {
	args := m.Called(serviceArea)
	return args.Error(0)
}

func (m *RabbitMQPublisher) UpdateServiceArea(ctx context.Context, serviceArea domain.ServiceArea) error {
	args := m.Called(serviceArea)
	return args.Error(0)
}
