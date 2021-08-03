package entities

type CityDocument struct {
	ID           string    `json:"_id"`
	Rev          string    `json:"_rev"`
	Key          string    `json:"_key"`
	Nom          string    `json:"nom"`
	CodesPostaux []string  `json:"codesPostaux"`
	Point        []float64 `json:"point"`
}
