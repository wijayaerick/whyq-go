package whyq

import (
	"context"
	"io"
	"net/http"
)

func (c *Client) Logout(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/logout", nil)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}
