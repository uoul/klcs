package main

import (
	"fmt"

	"github.com/uoul/go-common/auth"
	"github.com/uoul/go-common/config"
	"github.com/uoul/go-common/db"
	"github.com/uoul/go-common/log"
	"github.com/uoul/klcs/backend/core/api"
	"github.com/uoul/klcs/backend/core/domain"
	"github.com/uoul/klcs/backend/core/logic"
	"github.com/uoul/klcs/backend/core/services"
)

func main() {
	cp := config.NewEnvVarProvider()
	cf := setupDbConnection(cp)
	logger := log.NewConsoleLogger(
		log.StringToLogLevel(cp.StringOrDefault("KLCS_LOG_LVL", "INFO"), log.INFO),
	)
	authenticator := auth.NewKeyCloakAuthenticator[*domain.OidcUser](cp.StringOrDefault("KLCS_JWKS_URI", ""))
	printService := services.NewPrintService()
	logic := logic.NewLogic(cf, logger, printService)
	api := api.NewApi(logic, authenticator, printService)
	api.Run(cp.UInt16OrDefault("KLCS_HTTP_PORT", 8080))
}

func setupDbConnection(cp config.IConfigProvider) db.IConnectionFactory {
	return db.NewConnectionFactory(
		fmt.Sprintf(
			"host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
			cp.StringOrDefault("KLCS_DB_HOST", "localhost"),
			cp.Int16OrDefault("KLCS_DB_PORT", 5432),
			cp.StringOrDefault("KLCS_DB_USER", ""),
			cp.StringOrDefault("KLCS_DB_PW", ""),
			cp.StringOrDefault("KLCS_DB_NAME", "postgres"),
			cp.StringOrDefault("KLCS_DB_SSL", "verify-full"),
		),
		"postgres",
	)
}
