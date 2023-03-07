package main

import (
	"context"
	"fmt"
	"log"
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

	_, err := apiClient.Login(ctx, qvapay.LoginRequest{
		Email:    os.Getenv("QVAPAY_USER"),
		Password: os.Getenv("QVAPAY_PASSWORD"),
	})

	if err != nil {
		log.Fatal(err)
	}

	req := qvapay.EditMeRequest{
		Name:     "Wilson",
		Lastname: "Bolaños",
		Username: "wilson.bolanossoto",
		// Email:    "wilson.bolanossoto@gmail.com",
		// Bio:      "This is a Sample Bio Message",
	}
	// Get Me endpoints
	me, err := apiClient.EditMe(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(me)
}
