package entities

type ArangoResult struct {
	Key string `json:"_key"`
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}
