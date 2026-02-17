package config

import (
	"log/slog"
	"strings"
)

type AppConfig struct {
	LogLvl   string `default:"INFO"`
	KlcsHost string
	TimeZone string `default:"Europe/Vienna"`
	Printer  struct {
		Id      string
		UsbAddr string
		NetAddr string
	}
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
