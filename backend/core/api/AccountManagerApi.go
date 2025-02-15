package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/api/dto"
	"github.com/uoul/klcs/backend/core/domain"

	appError "github.com/uoul/klcs/backend/core/error"
)

func (e *Api) getAccounts(ctx *gin.Context) {
	accounts, err := e.logic.GetAllAccounts(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

func (e *Api) getAccountDetails(ctx *gin.Context) {
	accountId := ctx.Param("accountId")
	accountDetails, err := e.logic.GetAccountDetails(ctx, accountId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, accountDetails)
}

func (e *Api) createAccount(ctx *gin.Context) {
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

func (e *Api) updateAccount(ctx *gin.Context) {
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

func (e *Api) closeAccount(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
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

func (e *Api) postToAccount(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
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
