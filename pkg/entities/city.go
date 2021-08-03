package entities

type CityCollection struct {
	Key          string    `json:"_key,omitempty"`
	ID           string    `json:"_id,omitempty"`
	Rev          string    `json:"_rev,omitempty"`
	Nom          string    `json:"nom"`
	CodesPostaux []string  `json:"codesPostaux"`
	Point        []float64 `json:"point"`
}

type CitySearchBody struct {
	Nom string `json:"nom"`
}
