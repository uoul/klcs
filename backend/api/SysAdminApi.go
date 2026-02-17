package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/domain"
	appError "github.com/uoul/klcs/backend/core/error"
)

func (e *Api) createShop(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	var body domain.Shop
	if err := ctx.BindJSON(&body); err != nil {
		ctx.Error(appError.NewErrValidation("failed to parse shop - %v", err))
		return
	}
	shop, err := e.logic.CreateShop(
		ctx,
		user.GetUsername(),
		body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, shop)
}

func (e *Api) getShops(ctx *gin.Context) {
	shops, err := e.logic.GetShops(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, shops)
}

func (e *Api) updateShop(ctx *gin.Context) {
	var body domain.Shop
	if err := ctx.BindJSON(&body); err != nil {
		ctx.Error(appError.NewErrValidation("failed to parse shop - %v", err))
		return
	}
	if ctx.Param("shopId") != body.Id {
		ctx.Error(appError.NewErrValidation("shopId of uri does not match id from body (%s != %s)", ctx.Param("shopId"), body.Id))
		return
	}
	shop, err := e.logic.UpdateShop(
		ctx,
		body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, shop)
}

func (e *Api) deleteShop(ctx *gin.Context) {
	err := e.logic.DeleteShop(ctx, ctx.Param("shopId"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
