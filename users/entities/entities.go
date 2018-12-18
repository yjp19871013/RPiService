package entities

type CreateTokenRequest struct {
	Email    string `json:"email" binding:"required,email_validator"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
