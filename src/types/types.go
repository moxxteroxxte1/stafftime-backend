package types 

import "time"

type LoginRequest struct {
	Username string `json:"userName"`
	Password string `json:"password"`
	KeepLoggedIn bool `json:"keep"`
}

type LoginResponse struct {
	Token string `json:"token"`
	UserID uint `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}
