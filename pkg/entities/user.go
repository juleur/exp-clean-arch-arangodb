package entities

type AuthUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUserBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
