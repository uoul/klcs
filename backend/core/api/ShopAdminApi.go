package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/api/dto"
	"github.com/uoul/klcs/backend/core/domain"
	appError "github.com/uoul/klcs/backend/core/error"
)

func (e *Api) getArticlesForShop(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shopId := ctx.Param("shopId")
	articles, err := e.logic.GetArticlesForShop(ctx, user.GetUsername(), shopId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, articles)
}

func (e *Api) createArticle(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shopId := ctx.Param("shopId")
	article := domain.ArticleDetails{}
	err = ctx.BindJSON(&article)
	if err != nil {
		ctx.Error(appError.NewErrInvalidInput("failed to parse article - %v", err))
		return
	}
	a, err := e.logic.CreateArticle(
		ctx,
		user.GetUsername(),
		shopId,
		&article,
	)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, a)
}

func (e *Api) getArticle(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	articleId := ctx.Param("articleId")
	article, err := e.logic.GetArticle(ctx, user.GetUsername(), articleId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, article)
}

func (e *Api) updateArticle(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	articleId := ctx.Param("articleId")
	article := domain.ArticleDetails{}
	err = ctx.BindJSON(&article)
	if err != nil {
		ctx.Error(appError.NewErrInvalidInput("failed to parse article - %v", err))
		return
	}
	if articleId != article.Id {
		ctx.Error(appError.NewErrValidation("article id's does not match (%s != %s)", articleId, article.Id))
		return
	}
	a, err := e.logic.UpdateArticle(ctx, user.GetUsername(), &article)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, a)
}

func (e *Api) deleteArticle(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	articleId := ctx.Param("articleId")
	err = e.logic.DeleteArticle(ctx, user.GetUsername(), articleId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (e *Api) getPrintersForShop(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shopId := ctx.Param("shopId")
	printers, err := e.logic.GetPrintersForShop(ctx, user.GetUsername(), shopId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, printers)
}

func (e *Api) createPrinter(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shopId := ctx.Param("shopId")
	printer := domain.Printer{}
	err = ctx.BindJSON(&printer)
	if err != nil {
		ctx.Error(appError.NewErrInvalidInput("failed to parse printer - %v", err))
		return
	}
	p, err := e.logic.CreatePrinter(ctx, user.GetUsername(), shopId, &printer)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, p)
}

func (e *Api) deletePrinter(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	printerId := ctx.Param("printerId")
	err = e.logic.DeletePrinter(ctx, user.GetUsername(), printerId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (e *Api) getUsers(ctx *gin.Context) {
	users, err := e.logic.GetUsers(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (e *Api) getShopUsers(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shopId := ctx.Param("shopId")
	users, err := e.logic.GetShopUsers(ctx, user.GetUsername(), shopId)
	if err != nil {
		ctx.Error(err)
		return
	}
	result := []dto.ShopUserDto{}
	for u, roles := range users {
		result = append(result, dto.ShopUserDto{
			Id:        u.Id,
			Username:  u.Username,
			Name:      u.Name,
			ShopRoles: roles,
		})
	}
	ctx.JSON(http.StatusOK, result)
}

func (e *Api) addUserRoleForShop(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shopId := ctx.Param("shopId")
	userId := ctx.Param("userId")
	role := domain.Role{}
	err = ctx.BindJSON(&role)
	if err != nil {
		ctx.Error(appError.NewErrInvalidInput("failed to bind role - %v", err))
		return
	}
	err = e.logic.AddUserRole(ctx, user.GetUsername(), shopId, userId, role.Id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (e *Api) deleteUserRoleFromShop(ctx *gin.Context) {
	user, err := e.authenticator.GetIdentity(ctx.Request.Header)
	if err != nil {
		ctx.Error(appError.NewErrAuthentication("failed to get user identity - %s", err))
		return
	}
	shopId := ctx.Param("shopId")
	userId := ctx.Param("userId")
	roleId := ctx.Param("roleId")
	err = e.logic.DeleteUserRole(ctx, user.GetUsername(), shopId, userId, roleId)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (e *Api) getRoles(ctx *gin.Context) {
	roles, err := e.logic.GetRoles(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, roles)
}
