package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/apperror"
	"github.com/uoul/klcs/backend/core/domain"
)

func (e *Api) reprint(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")
	if err := e.logic.Reprint(ctx, transactionId); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (e *Api) getHistoryForUser(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
	if err != nil {
		ctx.Error(apperror.NewErrUnauthorized(err, "failed get user identity"))
		return
	}
	lengthStr := ctx.DefaultQuery("length", "10")
	len, err := strconv.Atoi(lengthStr)
	if err != nil {
		ctx.Error(apperror.NewErrLengthMustBeNumber(err, "length(%s) parameter has to be a number", lengthStr))
		return
	}
	history, err := e.logic.GetHistory(ctx, user.GetUsername(), len)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, history)
}

func (e *Api) getShopsForUser(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
	if err != nil {
		ctx.Error(apperror.NewErrUnauthorized(err, "failed get user identity"))
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
	user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
	if err != nil {
		ctx.Error(apperror.NewErrUnauthorized(err, "failed get user identity"))
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
	user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
	if err != nil {
		ctx.Error(apperror.NewErrUnauthorized(err, "failed get user identity"))
		return
	}
	reqOrder := domain.Order{}
	err = ctx.BindJSON(&reqOrder)
	if err != nil {
		ctx.Error(apperror.NewErrInvalidDataFormat(err, "failed to parse order"))
		return
	}
	resOrder, err := e.logic.Checkout(ctx, user.GetUsername(), user.HasRole(e.oidc.Roles.NoPrint), reqOrder)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, resOrder)
}
