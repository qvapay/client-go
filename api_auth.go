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

func (c *apiClient) Login(ctx context.Context, payload LoginRequest) (*LoginResponse, error) {

	url := fmt.Sprintf("%s/%s", c.server, loginEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload.ToReader())
	if err != nil {
		return authUser, fmt.Errorf("failed to create HTTP request  for Login: %v", err)
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return authUser, fmt.Errorf("failed to execute HTTP request  for Login: %v", err)
	}
	defer DrainBody(res.Body)

	if err = json.NewDecoder(res.Body).Decode(&authUser); err != nil {
		return nil, fmt.Errorf("failed to create HTTP response  for Login: %v", err)
	}

	return authUser, nil
}

func (c *apiClient) Register(ctx context.Context, payload RegisterRequest) (RegisterResponse, error) {

	url := fmt.Sprintf("%s/%s", c.server, registerEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload.ToReader())
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to create HTTP request for Register: %v", err)
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to execute HTTP request for Register: %v", err)
	}
	defer DrainBody(res.Body)

	var registeredUser RegisterResponse
	if err = json.NewDecoder(res.Body).Decode(&registeredUser); err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to create HTTP response for Register: %v", err)
	}

	return registeredUser, nil
}

func (c *apiClient) Logout(ctx context.Context) (LogoutResponse, error) {

	url := fmt.Sprintf("%s/%s", c.server, logoutEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return LogoutResponse{}, fmt.Errorf("failed to create HTTP request for Logout: %v", err)
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return LogoutResponse{}, fmt.Errorf("failed to execute HTTP request for Logout: %v", err)
	}
	defer DrainBody(res.Body)

	var logoutResp LogoutResponse
	if err = json.NewDecoder(res.Body).Decode(&logoutResp); err != nil {
		return LogoutResponse{}, fmt.Errorf("failed to create HTTP response for Logout: %v", err)
	}

	if res.StatusCode == http.StatusCreated {
		authUser.Clean()
	}

	return logoutResp, nil
}
