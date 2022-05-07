package dto

import (
	"service-area-service/internal/core/domain"
)

type serviceAreasResponse struct {
	ID         int    `json:"id"`
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}

func createServiceAreasResponse(serviceArea domain.ServiceArea) serviceAreasResponse {
	return serviceAreasResponse{
		ID:         serviceArea.ID,
		Identifier: serviceArea.Identifier,
		Name:       serviceArea.Name,
	}
}

type ServiceAreaListResponse []*serviceAreasResponse

func CreateServiceAreaListResponse(serviceAreas []domain.ServiceArea) ServiceAreaListResponse {
	response := ServiceAreaListResponse{}
	for _, s := range serviceAreas {
		serviceArea := createServiceAreasResponse(s)
		response = append(response, &serviceArea)
	}
	return response
}
