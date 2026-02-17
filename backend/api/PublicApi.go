package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/config"
)

func (a *Api) getUiSettings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, struct {
		Version string
		Oidc    config.OidcConfig
	}{
		Version: a.version,
		Oidc:    a.oidc,
	})
}
