package whyq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (c *Client) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	var loginResp LoginResponse

	form := url.Values{}
	form.Set("loginUser", email)
	form.Set("loginPass", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/ajax/login", strings.NewReader(form.Encode()))
	if err != nil {
		return loginResp, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := c.c.Do(req)
	if err != nil {
		return loginResp, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return loginResp, err
	}

	return loginResp, nil
}
