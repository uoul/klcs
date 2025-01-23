package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uoul/klcs/backend/oos-core/domain"
	appError "github.com/uoul/klcs/backend/oos-core/error"
)

func (e *ApiEnv) updateUserByOidcData() func(*gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentity(ctx.Request.Header)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		_, err = e.logic.UpdateUser(ctx, &domain.User{
			Username: user.UserName,
			Name:     user.Name,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Next()
	}
}

func (e *ApiEnv) checkUserLoggedIn() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentity(ctx.Request.Header)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		ctx.Set("oidcIdentity", *user)
		ctx.Next()
	}
}

func (e *ApiEnv) checkOidcRole(role string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		user, err := e.authenticator.GetIdentity(ctx.Request.Header)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if !user.HasRole(role) {
			ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("user %s does not have necessary role %s", user.GetUsername(), role))
			return
		}
		ctx.Next()
	}
}

func (e *ApiEnv) useCors() func(ctx *gin.Context) {
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

func (e *ApiEnv) errorTranslation() func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()
		if !ctx.IsAborted() {
			for _, err := range ctx.Errors {
				switch err.Err.(type) {
				case *appError.PermissionError:
					ctx.Status(http.StatusForbidden)
				case *appError.ValidationError:
					ctx.Status(http.StatusBadRequest)
				case *appError.NotFoundError:
					ctx.Status(http.StatusNotFound)
				default:
					ctx.Status(http.StatusInternalServerError)
				}
			}
		}
	}
}
