package dto

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
}

type User struct {
	ID       int
	Username string
	Password string
}
