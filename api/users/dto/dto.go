package dto

type (
	CreateTokenRequest struct {
		Email    string `json:"email" binding:"required,email_validator"`
		Password string `json:"password" binding:"required"`
	}

	TokenResponse struct {
		Token string   `json:"token"`
		Roles []string `json:"roles"`
	}

	RegisterRequest struct {
		Email        string `json:"email" binding:"required,email_validator"`
		Password1    string `json:"password1" binding:"required"`
		Password2    string `json:"password2" binding:"required"`
		ValidateCode string `json:"validateCode" binding:"required"`
	}

	ValidateCodeRequest struct {
		Email string `json:"email" binding:"required,email_validator"`
	}

	GetRolesResponse struct {
		Roles []string `json:"roles" binding:"required"`
	}
)
