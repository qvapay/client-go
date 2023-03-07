package qvapay

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginErrors(t *testing.T) {
	mockResp := `
	{
		"accessToken":"387003",
		"token_type":"Bearer",
		"me":{
			"uuid":"sd-sd-s-sd-sd",
			"username":"user@gmail.com",
			"name":"QvaPayUser",
			"lastname":"Paymenton"
		}
	}
	`

	table := []struct {
		Name       string
		Request    LoginRequest
		Response   string
		StatusCode int
		Want       error
	}{
		{Name: "Successful Login Request", StatusCode: http.StatusOK, Request: LoginRequest{Email: "mock@gmail.com", Password: "1q2w3e4r5t"}, Want: nil, Response: mockResp},
		{Name: "Empty Login Request", StatusCode: http.StatusOK, Request: LoginRequest{Email: "", Password: ""}, Want: ErrCreateReq, Response: mockResp},
		{Name: "Empty Password at Login", StatusCode: http.StatusOK, Request: LoginRequest{Email: "some@gm.com", Password: ""}, Want: ErrCreateReq, Response: mockResp},
		{Name: "Empty Email at Login Request", StatusCode: http.StatusOK, Request: LoginRequest{Email: "", Password: "1q2w3e4r5t"}, Want: ErrCreateReq, Response: mockResp},
		{Name: "Unmarshable Response", StatusCode: http.StatusOK, Request: LoginRequest{Email: "some@hm.com", Password: "1q2w3e4r5t"}, Want: ErrCreateRes, Response: "{" + mockResp},
		{Name: "Unmarshable Response", StatusCode: http.StatusUnprocessableEntity, Request: LoginRequest{Email: "some@hm.com", Password: "1q2w3e4r5t"}, Want: ErrUnsuccessfulRes, Response: mockResp},
	}

	for _, tt := range table {

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.StatusCode)
			fmt.Fprintf(w, tt.Response)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		_, err := api.Login(context.TODO(), tt.Request)

		assert.Equal(t, tt.Want, err)
	}
}

func TestLoginResponse(t *testing.T) {
	mockResp := `
	{
		"accessToken":"387003",
		"token_type":"Bearer",
		"me":{
			"uuid":"sd-sd-s-sd-sd",
			"username":"user@gmail.com",
			"name":"QvaPayUser",
			"lastname":"Paymenton"
		}
	}
	`

	table := []struct {
		Name    string
		Request LoginRequest
		Want    string
	}{
		{Name: "Successful Login Response", Request: LoginRequest{Email: "mock@gmail.com", Password: "1q2w3e4r5t"}, Want: "QvaPayUser"},
	}

	for _, tt := range table {

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, mockResp)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		login, err := api.Login(context.TODO(), tt.Request)
		if err != nil {
			t.Errorf("Unexpected error: %v", err.Error())
		}

		assert.Equal(t, tt.Want, login.Me.Name)
	}
}

func TestLogout(t *testing.T) {
	t.Run("sad path - User wasn't Logged in Before", func(t *testing.T) {
		data := `
	{
		"message":"successfully logout"
	}
	`
		authUser = nil

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, data)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		logout, err := api.Logout(context.TODO())
		if err != nil {
			assert.Equal(t, ErrCreateReq, err)
		}

		assert.Empty(t, logout)
	})

	t.Run("happy path - test http server successfully logout", func(t *testing.T) {
		data := `
	{
		"message":"successfully logout"
	}
	`
		authUser = &LoginResponse{
			AccessToken: "1q2w3e4r5t",
		}

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, data)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		logout, err := api.Logout(context.TODO())
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Equal(t, nil, err)
		assert.NotEmpty(t, logout)
		assert.Equal(t, "", authUser.AccessToken)
	})

	t.Run("sad path - Unmarshable response Body", func(t *testing.T) {
		data := `
	{{
		"message":"successfully logout"
	}
	`
		authUser = &LoginResponse{
			AccessToken: "1q2w3e4r5t",
		}

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, data)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		logout, err := api.Logout(context.TODO())
		if err != nil {
			assert.Equal(t, ErrCreateRes, err)
		}

		assert.Empty(t, logout)
		assert.NotEqual(t, "", authUser.AccessToken)
	})

}
