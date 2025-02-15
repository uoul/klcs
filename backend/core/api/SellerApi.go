package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/domain"
	appError "github.com/uoul/klcs/backend/core/error"
)

func (e *Api) getShopsForUser(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shops, err := e.logic.GetShopsForUser(ctx, user.GetUsername())
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, shops)
}

func (e *Api) getShopdetailsForUser(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
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

func (e *Api) checkout(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	reqOrder := domain.Order{}
	err = ctx.BindJSON(&reqOrder)
	if err != nil {
		ctx.Error(appError.NewErrInvalidInput("failed to parse order - %v", err))
		return
	}
	resOrder, err := e.logic.Checkout(ctx, user.GetUsername(), &reqOrder)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, resOrder)
}
