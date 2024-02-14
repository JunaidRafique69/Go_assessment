package models

// structs for authentication process

type Credentials struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type RefreshToken struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
