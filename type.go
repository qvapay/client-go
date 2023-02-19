package qvapay

// URL_BASE = https://qvapay.com/api

const (
	loginEndpoint    = "auth/login"
	registerEndpoint = "auth/register"
	logoutEndpoint   = "auth/logout"
)

var authUser *LoginResponse // Used as authentication state

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"token_type"`
	Me          User   `json:"me"`
}

func (l *LoginResponse) Clean() {
	l.AccessToken = ""
	l.TokenType = ""
	l.Me = User{}
}

type RegisterResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"token_type"`
	Me          User   `json:"me"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type User struct {
	UUID             string `json:"uuid"`
	Username         string `json:"username"`
	Name             string `json:"name"`
	Lastname         string `json:"lastname"`
	Bio              string `json:"bio"`
	ProfilePhotoPath string `json:"profile_photo_path"`
	Balance          int64  `json:"balance"`
	CompleteName     string `json:"complete_name"`
	NameVerified     string `json:"name_verified"`
	ProfilePhotoURL  string `json:"profile_photo_url"`
	AverageRating    string `json:"average_rating"`
}
