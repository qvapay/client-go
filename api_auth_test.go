package qvapay

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Run("happy path - test http server successfully", func(t *testing.T) {
		want := "QvaPayUser"
		data := `
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

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, data)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		login, err := api.Login(context.TODO(), LoginRequest{
			Email:    "user@gmail.com",
			Password: "CffdKB73iTtzNJN!",
		})
		if err != nil {
			t.Fatalf(err.Error())
		}
		assert.Equal(t, want, login.Me.Name)
	})

	t.Run("sad path - Empty Login payload values", func(t *testing.T) {
		data := `
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

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, data)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		_, err := api.Login(context.TODO(), LoginRequest{
			Email:    "",
			Password: "",
		})

		assert.Equal(t, ErrCreateReq, err)
	})

	t.Run("sad path - Wrong Credentials values", func(t *testing.T) {
		data := `
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

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, data)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		_, err := api.Login(context.TODO(), LoginRequest{
			Email:    "bademail@gial.com",
			Password: "1q2w3e4r5t",
		})

		assert.Equal(t, ErrUnsuccessfulRes, err)
	})

	t.Run("sad path - Unmarshable response Body", func(t *testing.T) {
		data := `
	{
		{"accessToken":"387003",
		"token_type":"Bearer",
		"me":{
			"uuid":"sd-sd-s-sd-sd",
			"username":"user@gmail.com",
			"name":"QvaPayUser",
			"lastname":"Paymenton"
		}
	}
	`

		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, data)
		}))

		defer svr.Close()
		api := NewAPIClient(APIClientOptions{
			Server: svr.URL,
		})

		_, err := api.Login(context.TODO(), LoginRequest{
			Email:    "user@gmail.com",
			Password: "1q2w3e4r5t",
		})

		assert.Equal(t, ErrCreateRes, err)
	})
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
