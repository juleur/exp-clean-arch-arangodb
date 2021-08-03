package entities

import (
	"time"
)

type UpdateActivityUser struct {
	Key      string    `json:"_key"`
	Datetime time.Time `json:"datetime"`
}
