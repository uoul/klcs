package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/config"
)

func (a *Api) getUiSettings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, struct {
		Version    string
		Oidc       config.OidcConfig
		UiSettings config.UiConfig
	}{
		Version:    a.version,
		Oidc:       a.oidc,
		UiSettings: a.uiConfig,
	})
}
