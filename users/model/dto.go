package model

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	Result string `json:"result"`
}
