package whyq

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Order struct {
	Date         time.Time
	Item         string
	OptionalItem string
}

func (c *Client) Orders(ctx context.Context) ([]Order, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/user/orders", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orders []Order

	tk := html.NewTokenizer(resp.Body)
	for {
		switch tk.Next() {
		case html.ErrorToken:
			err := tk.Err()
			if errors.Is(err, io.EOF) {
				return orders, nil
			}
			return nil, tk.Err()
		case html.StartTagToken:
			tkn := tk.Token()
			if tkn.Data != "table" {
				continue
			}

			attrs := nodeAttrsToMap(tkn.Attr)
			classVal, ok := attrs["class"]
			if !ok {
				continue
			}
			if strings.Contains(classVal, "table_headingcolor") && !strings.Contains(classVal, "table-striped") {
				order, ok, err := c.parseOrder(tk)
				if err != nil {
					return nil, err
				}
				if ok {
					orders = append(orders, order)
				}
			}
		}
	}
}

func (c *Client) parseOrder(tk *html.Tokenizer) (order Order, ok bool, err error) {
	var (
		foundDate, foundItemName bool
		isDate, isItemName       bool
	)

	for {
		switch tk.Next() {
		case html.ErrorToken:
			return order, ok, tk.Err()
		case html.EndTagToken:
			if tk.Token().Data != "table" {
				continue
			}
			return order, foundDate && foundItemName, nil
		case html.StartTagToken:
			tkn := tk.Token()
			if tkn.Data != "td" {
				continue
			}
			attrs := nodeAttrsToMap(tkn.Attr)
			if attrs["scope"] != "row" {
				continue
			}
			dataLabel, ok := attrs["data-label"]
			if !ok {
				continue
			}
			if strings.EqualFold(dataLabel, "Delivery Date") {
				isDate = true
				continue
			}
			if strings.EqualFold(dataLabel, "Item Name") {
				isItemName = true
				continue
			}
		case html.TextToken:
			tkn := tk.Token()
			if isDate {
				isDate = false
				foundDate = true
				date, err := time.Parse(timeLayout, strings.TrimSpace(tkn.Data))
				if err != nil {
					return order, ok, err
				}
				order.Date = date
				continue
			}
			if isItemName {
				isItemName = false
				foundItemName = true
				order.Item = strings.TrimSpace(tkn.Data)
				optionalItem, err := c.parseOptionalItem(tk)
				if err != nil {
					return order, ok, err
				}
				order.OptionalItem = optionalItem
			}
		}
	}
}

func (c *Client) parseOptionalItem(tk *html.Tokenizer) (item string, err error) {
	for {
		switch tk.Next() {
		case html.ErrorToken:
			return item, tk.Err()
		case html.EndTagToken:
			if tk.Token().Data != "td" {
				continue
			}
			return item, nil
		case html.TextToken:
			tkn := tk.Token()
			if strings.EqualFold(tkn.Data, "optional") {
				continue
			}
			item = strings.TrimPrefix(tkn.Data, ": ")
			item = strings.TrimSpace(item)
		}
	}
}

const timeLayout = "January 02, 2006"

func nodeAttrsToMap(attrs []html.Attribute) map[string]string {
	m := make(map[string]string)
	for _, attr := range attrs {
		m[attr.Key] = attr.Val
	}
	return m
}
