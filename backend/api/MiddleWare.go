package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/domain"
	appError "github.com/uoul/klcs/backend/core/error"
)

const (
	AUTH_HEADER = "Authorization"
)

func (e *Api) updateUserByOidcData() func(*gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				NewErrorResponse(appError.NewErrAuthentication("failed to get user identity - %v", err)),
			)
			return
		}
		_, err = e.logic.UpdateUser(ctx, domain.User{
			Username: user.UserName,
			Name:     user.Name,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}
		ctx.Next()
	}
}

func (e *Api) checkUserLoggedIn() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				NewErrorResponse(appError.NewErrAuthentication("failed to get user identity - %v", err)),
			)
			return
		}
		ctx.Next()
	}
}

func (e *Api) checkOidcRole(role string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				NewErrorResponse(appError.NewErrAuthentication("failed to get user identity - %v", err)),
			)
			return
		}
		if !user.HasRole(role) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, appError.NewErrForbidden("user %s does not have necessary role %s", user.GetUsername(), role))
			return
		}
		ctx.Next()
	}
}

func (e *Api) errorTranslation() func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()
		if !ctx.IsAborted() {
			for _, err := range ctx.Errors {
				resp := NewErrorResponse(err.Err)
				switch err.Err.(type) {
				case appError.ErrForbidden:
					ctx.JSON(http.StatusForbidden, resp)
				case appError.ErrValidation:
					ctx.JSON(http.StatusBadRequest, resp)
				case appError.ErrNotFound:
					ctx.JSON(http.StatusNotFound, resp)
				case appError.ErrAuthentication:
					ctx.JSON(http.StatusUnauthorized, resp)
				default:
					ctx.JSON(http.StatusInternalServerError, resp)
				}
			}
		}
	}
}
