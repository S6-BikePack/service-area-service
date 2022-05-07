package dto

import "service-area-service/internal/core/domain"

type ServiceAreaResponse struct {
	ID         int         `json:"id"`
	Identifier string      `json:"identifier"`
	Name       string      `json:"name"`
	Area       domain.Area `json:"area"`
}

func CreateServiceAreaResponse(serviceArea domain.ServiceArea) ServiceAreaResponse {
	return ServiceAreaResponse{
		ID:         serviceArea.ID,
		Identifier: serviceArea.Identifier,
		Name:       serviceArea.Name,
		Area:       serviceArea.Area,
	}
}
