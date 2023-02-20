package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	qvapay "github.com/qvapay/client-go"
)

func main() {
	ctx := context.Background()
	apiClient := qvapay.NewAPIClient(qvapay.APIClientOptions{
		HttpClient: http.DefaultClient,
		Debug:      os.Stdout,
		Server:     "https://qvapay.com",
	})

	resultLogin, err := apiClient.Login(ctx, qvapay.LoginRequest{
		Email:    os.Getenv("QVAPAY_USER"),
		Password: os.Getenv("QVAPAY_PASSWORD"),
	})

	// defaults errors
	if err != nil {
		//E.g: How you can use custom errors
		if errors.Is(err, qvapay.ErrExecuteReq) {
			fmt.Println(err)
		}
		fmt.Println(err)
	}
	fmt.Println("Result Before Logout:", resultLogin)

	resultLogout, err := apiClient.Logout(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Result:", resultLogout)
	fmt.Println("Result After Logout:", resultLogin)
}
