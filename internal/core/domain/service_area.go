package domain

type ServiceArea struct {
	ID            int
	Identifier    string
	Name          string
	RiderCoverage int
	Area          Area
}

func NewServiceArea(id int, identifier string, name string, area Area) ServiceArea {
	return ServiceArea{
		ID:         id,
		Identifier: identifier,
		Name:       name,
		Area:       area,
	}
}
