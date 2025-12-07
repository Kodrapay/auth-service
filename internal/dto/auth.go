package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int64   `json:"expires_in"`
	MerchantID   *string `json:"merchant_id,omitempty"`
	Role         string  `json:"role,omitempty"`
	SessionID    string  `json:"session_id,omitempty"`
	Email        string  `json:"email,omitempty"`
}

type RegisterRequest struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	Name       string  `json:"name"`
	MerchantID *string `json:"merchant_id,omitempty"`
}

type RegisterResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	MerchantID   string `json:"merchant_id,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type ValidateSessionRequest struct {
	SessionID string `json:"session_id"`
}

type ValidateSessionResponse struct {
	Valid      bool    `json:"valid"`
	UserID     string  `json:"user_id,omitempty"`
	Role       string  `json:"role,omitempty"`
	MerchantID string  `json:"merchant_id,omitempty"`
	Email      string  `json:"email,omitempty"`
}
