package model

type AuthCredentials struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SigninCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
