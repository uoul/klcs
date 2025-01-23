package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/oos-core/domain"
	appError "github.com/uoul/klcs/backend/oos-core/error"
)

func (e *ApiEnv) createShop(ctx *gin.Context) {
	user, err := getFromRequestCtx[domain.OidcUser](ctx, "oidcIdentity")
	if err != nil {
		ctx.Error(err)
		return
	}
	var body domain.Shop
	if err := ctx.BindJSON(&body); err != nil {
		ctx.Error(appError.NewValidationError(err))
		return
	}
	shop, err := e.logic.CreateShop(
		ctx,
		user.GetUsername(),
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, shop)
}

func (e *ApiEnv) getShops(ctx *gin.Context) {
	shops, err := e.logic.GetShops(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, shops)
}

func (e *ApiEnv) updateShop(ctx *gin.Context) {
	var body domain.Shop
	if err := ctx.BindJSON(&body); err != nil {
		ctx.Error(appError.NewValidationError(err))
		return
	}
	if ctx.Param("shopId") != body.Id {
		ctx.Error(appError.NewValidationError(fmt.Errorf("shopIds does not match in request")))
		return
	}
	shop, err := e.logic.UpdateShop(
		ctx,
		&body,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, shop)
}

func (e *ApiEnv) deleteShop(ctx *gin.Context) {
	err := e.logic.DeleteShop(ctx, ctx.Param("shopId"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
