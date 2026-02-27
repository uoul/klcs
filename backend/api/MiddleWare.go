package api

import (
	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/apperror"
	"github.com/uoul/klcs/backend/core/domain"
)

const (
	AUTH_HEADER = "Authorization"
)

func (e *Api) updateUserByOidcData() func(*gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
		if err != nil {
			ctx.Error(apperror.NewErrUnauthorized(err, "failed to get user identity"))
			ctx.Abort()
			return
		}
		_, err = e.logic.UpdateUser(ctx, domain.User{
			Username: user.UserName,
			Name:     user.Name,
		})
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (e *Api) checkUserLoggedIn() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
		if err != nil {
			ctx.Error(apperror.NewErrUnauthorized(err, "failed to get user identity"))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (e *Api) checkOidcRole(role string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentityFromHeader(ctx.Request.Header, AUTH_HEADER)
		if err != nil {
			ctx.Error(apperror.NewErrUnauthorized(err, "failed to get user identity"))
			ctx.Abort()
			return
		}
		if !user.HasRole(role) {
			ctx.Error(apperror.NewErrMissingOidcRole(nil, "user %s does not have necessary role %s", user.GetUsername(), role))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func (e *Api) errorTranslation() func(*gin.Context) {
	return func(ctx *gin.Context) {
		// Execute all other handlers first
		ctx.Next()
		// Only handle requests with errors
		if len(ctx.Errors) <= 0 {
			return // do nothing
		}
		// Get last error
		err := ctx.Errors.Last()
		// Check if error is IAppError
		appErr, isAppError := err.Err.(apperror.IAppError)
		if isAppError {
			ctx.JSON(appErr.HttpStatus(), e.NewErrorResponse(appErr))
			return
		}
		// Otherwise return default error
		defaultErr := apperror.NewErrInternal(err, "Internal server error")
		ctx.JSON(defaultErr.HttpStatus(), e.NewErrorResponse(defaultErr))
	}
}
