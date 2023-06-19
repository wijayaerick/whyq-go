package whyq

import (
	"net/http"
)

type HttpDoer interface {
	Do(r *http.Request) (*http.Response, error)
}

type Client struct {
	c       HttpDoer
	baseURL string
}

func NewClient(c HttpDoer, baseURL string) *Client {
	return &Client{
		c:       c,
		baseURL: baseURL,
	}
}
