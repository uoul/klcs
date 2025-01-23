package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/oos-core/api/dto"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

func (e *ApiEnv) getAccountDetails(ctx *gin.Context) {
	accountId := ctx.Param("accountId")
	accountDetails, err := e.logic.GetAccountDetails(ctx, accountId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, accountDetails)
}

func (e *ApiEnv) createAccount(ctx *gin.Context) {
	account := domain.Account{}
	err := ctx.BindJSON(&account)
	if err != nil {
		ctx.Error(err)
		return
	}
	result, err := e.logic.CreateAccount(ctx, &account)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (e *ApiEnv) updateAccount(ctx *gin.Context) {
	account := domain.Account{}
	err := ctx.BindJSON(&account)
	if err != nil {
		ctx.Error(err)
		return
	}
	result, err := e.logic.UpdateAccount(ctx, &account)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (e *ApiEnv) closeAccount(ctx *gin.Context) {
	user, err := getFromRequestCtx[domain.OidcUser](ctx, "oidcIdentity")
	if err != nil {
		ctx.Error(err)
		return
	}
	accountId := ctx.Param("accountId")
	account, err := e.logic.CloseAccount(ctx, user.GetUsername(), accountId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (e *ApiEnv) postToAccount(ctx *gin.Context) {
	user, err := getFromRequestCtx[domain.OidcUser](ctx, "oidcIdentity")
	if err != nil {
		ctx.Error(err)
		return
	}
	accountId := ctx.Param("accountId")
	accountBalanceUpdate := dto.AccountBalanceUpdateDto{}
	err = ctx.BindJSON(&accountBalanceUpdate)
	if err != nil {
		ctx.Error(err)
		return
	}
	accountDetails, err := e.logic.PostToAccount(ctx, user.GetUsername(), accountId, accountBalanceUpdate.Amount)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, accountDetails)
}
