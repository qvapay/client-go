package qvapay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type MeRAW map[string]any

type EditMeRequest struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Bio      string `json:"bio"`
	Logo     string `json:"logo"`
	KYC      int    `json:"kyc"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *EditMeRequest) ToReader() *strings.Reader {
	byte, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}
	return strings.NewReader(string(byte))
}

func (p *EditMeRequest) Validate() error {
	if reflect.ValueOf(p).IsZero() {
		return fmt.Errorf("invalid Payload trying to be passed for User edition")
	}

	return nil
}

func (c *apiClient) GetMeRAW(ctx context.Context) (MeRAW, error) {
	url := fmt.Sprintf("%s/%s", c.server, meEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return MeRAW{}, ErrCreateReq
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return MeRAW{}, ErrExecuteReq
	}
	defer DrainBody(res.Body)

	if res.StatusCode != http.StatusOK {
		return MeRAW{}, ErrUnsuccessfulRes
	}

	var result MeRAW

	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, ErrCreateRes
	}

	return result, nil
}

func (c *apiClient) EditMe(ctx context.Context, payload EditMeRequest) (User, error) {

	if err := payload.Validate(); err != nil {
		return User{}, ErrCreateReq
	}

	url := fmt.Sprintf("%s/%s", c.server, meEndpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, payload.ToReader())
	if err != nil {
		return User{}, ErrCreateReq
	}

	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return User{}, ErrExecuteReq
	}
	defer DrainBody(res.Body)

	if res.StatusCode != http.StatusCreated {
		return User{}, ErrUnsuccessfulRes
	}

	var result User

	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return User{}, ErrCreateRes
	}

	// Update the state with the new User data
	authUser.Me = result

	return result, nil
}
