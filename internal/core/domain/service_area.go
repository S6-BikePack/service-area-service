package domain

type ServiceArea struct {
	ID            int    `json:"id"`
	Identifier    string `json:"identifier"`
	Name          string `json:"name"`
	RiderCoverage int    `json:"riderCoverage"`
	Area          Area   `json:"area"`
}

func NewServiceArea(id int, identifier string, name string, area Area) ServiceArea {
	return ServiceArea{
		ID:         id,
		Identifier: identifier,
		Name:       name,
		Area:       area,
	}
}
