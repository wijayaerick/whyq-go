package internal

import (
	"context"
	"log/slog"
	"os"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	BaseURL  string `env:"WHYQ_URL,required"`
	Email    string `env:"WHYQ_EMAIL"`
	Password string `env:"WHYQ_PASSWORD"`
}

func LoadConfig(ctx context.Context) (Config, error) {
	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return c, err
	}
	return c, nil
}

func (c Config) Logger() *slog.Logger {
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelInfo,
		ReplaceAttr: nil,
	})
	return slog.New(textHandler)
}
