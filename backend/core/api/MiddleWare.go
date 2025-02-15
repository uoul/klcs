package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/core/domain"
	appError "github.com/uoul/klcs/backend/core/error"
)

func (e *Api) updateUserByOidcData() func(*gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentity(ctx.Request.Header)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				NewErrorResponse(appError.NewErrAuthentication("failed to get user identity - %v", err)),
			)
			return
		}
		_, err = e.logic.UpdateUser(ctx, &domain.User{
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
		_, err := e.authenticator.GetIdentity(ctx.Request.Header)
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
		user, err := e.authenticator.GetIdentity(ctx.Request.Header)
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

func (e *Api) useCors() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
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
