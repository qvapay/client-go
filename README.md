# client-go
GO client for the QvaPay API



## Examples

```go

package main

import (
	"context"
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
		Server:     "https://qvapay.com/api",
	})

	result, err := apiClient.Login(ctx, qvapay.LoginRequest{
		Email:    os.Getenv("QVAPAY_USER"),
		Password: os.Getenv("QVAPAY_PASSWORD"),
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

```

For more examples take a look in `examples` directory.