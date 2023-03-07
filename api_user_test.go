package qvapay

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMeErrors(t *testing.T) {
	mockResp := `
	{
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
