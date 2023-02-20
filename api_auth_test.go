package qvapay

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginOK(t *testing.T) {

	t.Run("happy path - test http server ", func(t *testing.T) {
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
}
