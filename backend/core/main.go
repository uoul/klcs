package main

import (
	"fmt"

	"github.com/uoul/go-common/auth/oidc"
	"github.com/uoul/go-common/config"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/api"
	"github.com/uoul/klcs/backend/oos-core/domain"
	"github.com/uoul/klcs/backend/oos-core/logic"
)

func main() {
	cp := config.NewEnvVarProvider()
	cf := setupDbConnection(cp)
	logic := logic.NewLogic(cf)
	authenticator := oidc.NewKeyCloakAuthenticator[*domain.OidcUser](cp.StringOrDefault("KLCS_JWKS_URI", ""))
	api := api.NewApiEnv(logic, authenticator)
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
			cp.StringOrDefault("KLCS_DB_SSL", "enabled"),
		),
		"postgres",
	)
}
