package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	AccountID   string `json:"accountId"`
}
