package services

import (
	"context"
	"service-area-service/config"
	"service-area-service/internal/core/domain"
	"service-area-service/internal/core/interfaces"
)

type serviceAreaService struct {
	config                *config.Config
	serviceAreaRepository interfaces.ServiceAreaRepository
	messagePublisher      interfaces.MessageBusPublisher
}

func NewServiceAreaService(config *config.Config, messagePublisher interfaces.MessageBusPublisher, serviceAreaRepository interfaces.ServiceAreaRepository) *serviceAreaService {
	return &serviceAreaService{
		config:                config,
		serviceAreaRepository: serviceAreaRepository,
		messagePublisher:      messagePublisher,
	}
}

func (srv *serviceAreaService) GetAll(ctx context.Context) ([]domain.ServiceArea, error) {
	return srv.serviceAreaRepository.GetAll(ctx)
}

func (srv *serviceAreaService) Get(ctx context.Context, id int) (domain.ServiceArea, error) {
	serviceArea, err := srv.serviceAreaRepository.Get(ctx, id)

	if err != nil {
		return domain.ServiceArea{}, err
	}

	return serviceArea, nil
}

func (srv *serviceAreaService) Create(ctx context.Context, id int, identifier, name string, area domain.Area) (domain.ServiceArea, error) {
	serviceArea := domain.NewServiceArea(id, identifier, name, area)

	serviceArea, err := srv.serviceAreaRepository.Save(ctx, serviceArea)

	if err != nil {
		return domain.ServiceArea{}, err
	}

	err = srv.messagePublisher.CreateServiceArea(ctx, serviceArea)

	return serviceArea, err
}

func (srv *serviceAreaService) Update(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error) {
	existing, err := srv.serviceAreaRepository.Get(ctx, serviceArea.ID)
	updated := existing

	if err != nil {
		return domain.ServiceArea{}, err
	}

	if serviceArea.Name != "" {
		updated.Name = serviceArea.Name
	}

	if serviceArea.Identifier != "" {
		updated.Identifier = serviceArea.Identifier
	}

	if serviceArea.Area.Type == "Polygon" {
		updated.Area = serviceArea.Area
	}

	serviceArea, err = srv.serviceAreaRepository.Update(ctx, updated)

	if err != nil {
		return existing, err
	}

	err = srv.messagePublisher.UpdateServiceArea(ctx, serviceArea)

	return serviceArea, err
}
