package dto

import (
	"service-area-service/internal/core/domain"
)

type BodyCreateServiceArea struct {
	ID         int         `json:"id"`
	Identifier string      `json:"identifier"`
	Name       string      `json:"name"`
	Area       domain.Area `json:"area"`
}
