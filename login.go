package whyq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type loginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (c *Client) Login(ctx context.Context, email, password string) error {
	var loginResp loginResponse

	form := url.Values{}
	form.Set("loginUser", email)
	form.Set("loginPass", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/ajax/login", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return err
	}

	if !strings.EqualFold(loginResp.Status, "success") {
		return fmt.Errorf("whyq: failed to login: %v", loginResp.Message)
	}

	return nil
}
