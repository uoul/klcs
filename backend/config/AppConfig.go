package config

import (
	"log/slog"
	"strings"
)

type AppConfig struct {
	LogLvl   string `default:"INFO"`
	Api      string `default:":80"`
	Debug    bool   `default:"false"`
	Cors     CorsConfig
	Oidc     OidcConfig
	Database DatabaseConfig `gonfig:"db"`
	Ui       UiConfig
}

func (c *AppConfig) SlogLvl() slog.Level {
	lvl := strings.ToUpper(c.LogLvl)
	switch lvl {
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	case "DEBUG":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}
