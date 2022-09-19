package internal

// LoginRequest request to login into the application.
type LoginRequest struct {
	UserOrEmail string `json:"user_or_email"`
	Password    string `json:"password"`
}

// LoginResponse response from logging in.
type LoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// BuildLoginResponse builds a response from logging in.
func BuildLoginResponse(username, email, token string) *LoginResponse {
	var logRes LoginResponse

	logRes.Username = username
	logRes.Email = email
	logRes.Token = token

	return &logRes
}
