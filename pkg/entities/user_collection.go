package entities

import "time"

type UserDoc struct {
	Key          string     `json:"_key,omitempty"`
	ID           string     `json:"_id,omitempty"`
	Rev          string     `json:"_rev,omitempty"`
	Pseudo       string     `json:"pseudo,omitempty"`
	Email        string     `json:"email,omitempty"`
	Hpwd         string     `json:"hpwd,omitempty"`
	RegisteredAt time.Time  `json:"registeredAt,omitempty"`
	LastSeenAt   *time.Time `json:"lastSeenAt,omitempty"`
}
