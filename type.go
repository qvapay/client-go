package qvapay

import "fmt"

// URL_BASE = https://qvapay.com/api

const (
	loginEndpoint    = "api/auth/login"
	registerEndpoint = "api/auth/register"
	logoutEndpoint   = "api/auth/logout"
)

// Customs Errors

var (
	ErrCreateReq       = fmt.Errorf("failed to create HTTP request")
	ErrExecuteReq      = fmt.Errorf("failed to execute HTTP request")
	ErrCreateRes       = fmt.Errorf("failed to create HTTP response")
	ErrUnsuccessfulRes = fmt.Errorf("failed retrieving a successful HTTP response due to a possible non well formed request Body")
)

// TODO: move to api_user.go
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
