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

## Tests

To run tests and check the coverage report, run the following commands:

``` sh
# Running simple tests
go test ./... -v

# Simple test coverage percentage
go test -cover

# Build test coverage profile
go test -coverprofile=cover.out

# Analyze coverage profile report in HTML
go tool cover -html=cover.out -o cover.html
```

For more examples take a look in `examples` directory.
