package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/oos-core/domain"
	appError "github.com/uoul/klcs/backend/oos-core/error"
)

func (e *ApiEnv) getShopsForUser(ctx *gin.Context) {
	user, err := getFromRequestCtx[domain.OidcUser](ctx, "oidcIdentity")
	if err != nil {
		ctx.Error(err)
		return
	}
	shops, err := e.logic.GetShopsForUser(ctx, user.GetUsername())
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, shops)
}

func (e *ApiEnv) getShopdetailsForUser(ctx *gin.Context) {
	user, err := getFromRequestCtx[domain.OidcUser](ctx, "oidcIdentity")
	if err != nil {
		ctx.Error(err)
		return
	}
	shopId := ctx.Param("shopId")
	shop, err := e.logic.GetShopDetailsForUser(ctx, user.GetUsername(), shopId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, shop)
}

func (e *ApiEnv) checkout(ctx *gin.Context) {
	user, err := getFromRequestCtx[domain.OidcUser](ctx, "oidcIdentity")
	if err != nil {
		ctx.Error(err)
		return
	}
	reqOrder := domain.Order{}
	err = ctx.BindJSON(&reqOrder)
	if err != nil {
		ctx.Error(appError.NewValidationError(err))
		return
	}
	resOrder, err := e.logic.Checkout(ctx, user.GetUsername(), &reqOrder)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, resOrder)
}
