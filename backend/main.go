package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/uoul/go-auth"
	gonfig "github.com/uoul/go-config"
	"github.com/uoul/klcs/backend/core/api"
	"github.com/uoul/klcs/backend/core/config"
	"github.com/uoul/klcs/backend/core/domain"
	"github.com/uoul/klcs/backend/core/logic"
	"github.com/uoul/klcs/backend/core/services"
)

const (
	VERSION = "{VERSION}"
)

func main() {
	// Get AppConfig
	cfg, err := gonfig.FromEnvironment[config.AppConfig]("KLCS")
	if err != nil {
		panic(err)
	}
	// Setup Logger
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     cfg.SlogLvl(),
	})
	slog.SetDefault(slog.New(logHandler))
	// Create database conneciton
	dbConn, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	))
	if err != nil {
		panic(err)
	}
	// Create Authenticator
	authenticator := auth.NewJwksAuthenticator[domain.OidcUser](cfg.Oidc.JwksUrl)
	// Run Services
	printService := services.NewPrintService()
	// Create Logic
	logic := logic.NewLogic(dbConn, printService)
	// Run Api
	if err := api.NewApi(
		VERSION,
		cfg.Cors,
		cfg.Oidc,
		logic,
		authenticator,
		printService,
		api.WithCorsOrigins(cfg.Cors.Origins),
		api.WithCorsHeaders(cfg.Cors.Headers),
		api.WithCorsMethods(cfg.Cors.Methods),
	).Run(cfg.Api); err != nil {
		slog.Error("failed to run api", slog.Any("error", err))
	}
}
