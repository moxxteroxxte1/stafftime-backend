package types

type LoginRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	KeepLoggedIn bool   `json:"keep"`
}

type LoginResponse struct {
	UserID  uint `json:"userId"`
	IsAdmin bool `json:"isAdmin"`
}
