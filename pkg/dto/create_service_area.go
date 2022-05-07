package dto

import (
	"service-area-service/internal/core/domain"
)

type BodyCreateServiceArea struct {
	ID         int
	Identifier string
	Name       string
	Area       domain.Area
}

type ResponseCreateServiceArea domain.ServiceArea

func BuildResponseCreate(model domain.ServiceArea) ResponseCreateServiceArea {
	return ResponseCreateServiceArea(model)
}
