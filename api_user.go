package qvapay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type MeRAW map[string]any

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
