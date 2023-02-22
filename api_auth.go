package qvapay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *LoginRequest) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"c_password"`
	Invite          string `json:"referer_username"`
}

func (p *RegisterRequest) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}

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

func (c *apiClient) Login(ctx context.Context, payload LoginRequest) (*LoginResponse, error) {

	if payload.Email == "" || payload.Password == "" {
		return nil, ErrCreateReq
	}

	url := fmt.Sprintf("%s/%s", c.server, loginEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload.ToReader())
	if err != nil {
		return authUser, ErrCreateReq
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return authUser, ErrExecuteReq
	}
	defer DrainBody(res.Body)

	if res.StatusCode != http.StatusOK {
		return authUser, ErrUnsuccessfulRes
	}

	if err = json.NewDecoder(res.Body).Decode(&authUser); err != nil {
		return nil, ErrCreateRes
	}

	return authUser, nil
}

func (c *apiClient) Register(ctx context.Context, payload RegisterRequest) (RegisterResponse, error) {

	url := fmt.Sprintf("%s/%s", c.server, registerEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload.ToReader())
	if err != nil {
		return RegisterResponse{}, ErrCreateReq
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return RegisterResponse{}, ErrExecuteReq
	}
	defer DrainBody(res.Body)

	var registeredUser RegisterResponse
	if err = json.NewDecoder(res.Body).Decode(&registeredUser); err != nil {
		return RegisterResponse{}, ErrCreateRes
	}

	return registeredUser, nil
}

func (c *apiClient) Logout(ctx context.Context) (LogoutResponse, error) {

	// Check that user must be logged in before logout
	if authUser == nil {
		return LogoutResponse{}, ErrCreateReq
	}

	if authUser.AccessToken == "" {
		return LogoutResponse{}, ErrCreateReq
	}

	url := fmt.Sprintf("%s/%s", c.server, logoutEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return LogoutResponse{}, ErrCreateReq
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return LogoutResponse{}, ErrExecuteReq
	}
	defer DrainBody(res.Body)

	var logoutResp LogoutResponse
	if err = json.NewDecoder(res.Body).Decode(&logoutResp); err != nil {
		return LogoutResponse{}, ErrCreateRes
	}

	if res.StatusCode != http.StatusCreated {
		return logoutResp, ErrUnsuccessfulRes
	}

	authUser.Clean()

	return logoutResp, nil
}
