package repositories

import (
	"context"
	"gorm.io/gorm"
	"service-area-service/internal/core/domain"
)

type ServiceAreaRepository struct {
	Connection *gorm.DB
}

func NewServiceAreaRepository(db *gorm.DB) (*ServiceAreaRepository, error) {
	_ = db.Exec(`CREATE EXTENSION IF NOT EXISTS postgis;`)
	err := db.AutoMigrate(&domain.ServiceArea{})

	if err != nil {
		return nil, err
	}

	database := ServiceAreaRepository{
		Connection: db,
	}

	return &database, nil
}

func (repository *ServiceAreaRepository) Get(ctx context.Context, id int) (domain.ServiceArea, error) {
	var serviceArea domain.ServiceArea

	repository.Connection.WithContext(ctx).First(&serviceArea, "id = ?", id)

	return serviceArea, nil
}

func (repository *ServiceAreaRepository) GetAll(ctx context.Context) ([]domain.ServiceArea, error) {
	var serviceArea []domain.ServiceArea

	repository.Connection.WithContext(ctx).Find(&serviceArea)

	return serviceArea, nil
}

func (repository *ServiceAreaRepository) Save(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error) {
	result := repository.Connection.WithContext(ctx).Create(&serviceArea)

	if result.Error != nil {
		return domain.ServiceArea{}, result.Error
	}

	return serviceArea, nil
}

func (repository *ServiceAreaRepository) Update(ctx context.Context, serviceArea domain.ServiceArea) (domain.ServiceArea, error) {
	result := repository.Connection.WithContext(ctx).Model(&serviceArea).Updates(serviceArea)

	if result.Error != nil {
		return domain.ServiceArea{}, result.Error
	}

	return serviceArea, nil
}
