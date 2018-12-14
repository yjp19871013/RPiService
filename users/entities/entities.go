package entities

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
