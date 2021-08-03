package entities

import "time"

type AlimentCollection struct {
	ImgUrl         string    `json:"imgUrl,omitempty"`
	Categorie      string    `json:"categorie"`
	Nom            string    `json:"nom"`
	Variete        string    `json:"variete"`
	SystemeEchange []int     `json:"systemeEchange"`
	Prix           float64   `json:"prix"`
	UniteMesure    int       `json:"uniteMesure"`
	Stock          int       `json:"stock"`
	AddedAt        time.Time `json:"addedAt"`
}

type RelUsersAliments struct {
	From string `json:"_from"`
	To   string `json:"_to"`
}
