package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/wijayaerick/whyq-go"
	"github.com/wijayaerick/whyq-go/internal"
	"golang.org/x/net/publicsuffix"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()

	conf, err := internal.LoadConfig(ctx)
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		return 1
	}
	logger := conf.Logger()

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		logger.ErrorCtx(ctx, "failed to create cookie jar", "err", err)
		return 1
	}

	client := whyq.NewClient(&http.Client{Jar: jar}, "https://www.whyq.sg")
	if err := client.Login(ctx, conf.Email, conf.Password); err != nil {
		logger.ErrorCtx(ctx, "failed to login", "err", err)
		return 1
	}

	defer func(ctx context.Context) {
		if err := client.Logout(ctx); err != nil {
			logger.ErrorCtx(ctx, "failed to logout", "err", err)
		}
	}(ctx)

	orders, err := client.Orders(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to get orders", "err", err)
		return 1
	}

	if len(orders) == 0 {
		fmt.Printf("No order.")
		return 0
	}
	for i, order := range orders {
		fmt.Printf("%2d. %v %v %v\n", i+1,
			ColorText(Green, order.Date.Weekday().String()[:3]),
			ColorText(Green, order.Date.Format(time.DateOnly)),
			ColorText(Cyan, order.Item))
		fmt.Printf("    Optional: %v\n", order.OptionalItem)
	}

	return 0
}

func ColorText(c uint8, s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}

const (
	Black uint8 = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)
