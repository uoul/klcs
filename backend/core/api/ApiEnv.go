package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/uoul/go-common/auth/iface"
	"github.com/uoul/klcs/backend/oos-core/domain"
	appError "github.com/uoul/klcs/backend/oos-core/error"
	"github.com/uoul/klcs/backend/oos-core/logic"
)

const (
	OIDC_ADMIN_ROLE           = "KLCS_ADMIN"
	OIDC_ACCOUNT_MANAGER_ROLE = "KLCS_ACCOUNT_MANAGER"
)

type ApiEnv struct {
	logic         logic.ILogic
	authenticator iface.IAuthenticator[*domain.OidcUser]
}

func (e *ApiEnv) Run(port uint16) {
	// Get new router
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		e.useCors(),
	)
	rootGroup := router.Group("/api/v1")
	// Setup global middleware
	rootGroup.Use(
		e.checkUserLoggedIn(),
		e.updateUserByOidcData(),
		e.errorTranslation(),
	)
	// Setup routergroups
	e.setupSysAdminRg(rootGroup, "/admin")
	e.setupUserRg(rootGroup, "")
	e.setupAccountManagerRg(rootGroup, "/accounts")
	// Run api
	router.Run(fmt.Sprintf(":%v", port))
}

func (e *ApiEnv) setupSysAdminRg(router *gin.RouterGroup, prefix string) *gin.RouterGroup {
	rg := router.Group(prefix)
	rg.Use(
		e.checkOidcRole(OIDC_ADMIN_ROLE),
	)
	// assign handlers
	rg.POST("/shops", e.createShop)
	rg.GET("/shops", e.getShops)
	rg.PATCH("/shops/:shopId", e.updateShop)
	rg.DELETE("/shops/:shopId", e.deleteShop)
	return rg
}

func (e *ApiEnv) setupAccountManagerRg(router *gin.RouterGroup, prefix string) {
	rg := router.Group(prefix)
	rg.Use(
		e.checkOidcRole(OIDC_ACCOUNT_MANAGER_ROLE),
	)
	rg.GET("/:accountId", e.getAccountDetails)
	rg.POST("", e.createAccount)
	rg.PATCH("/:accountId", e.updateAccount)
	rg.DELETE("/:accountId/balance", e.closeAccount)
	rg.POST("/:accountId/balance", e.postToAccount)
}

func (e *ApiEnv) setupUserRg(router *gin.RouterGroup, prefix string) *gin.RouterGroup {
	rg := router.Group(prefix)
	// seller api
	rg.GET("/shops", e.getShopsForUser)
	rg.GET("/shops/:shopId", e.getShopdetailsForUser)
	rg.POST("/orders", e.checkout)
	// shopadmin api
	rg.GET("/shops/:shopId/articles", e.getArticlesForShop)
	rg.POST("/shops/:shopId/articles", e.createArticle)
	rg.GET("/articles/:articleId", e.getArticle)
	rg.PATCH("/articles/:articleId", e.updateArticle)
	rg.DELETE("/articles/:articleId", e.deleteArticle)
	rg.GET("/shops/:shopId/printers", e.getPrintersForShop)
	rg.POST("/shops/:shopId/printers", e.createPrinter)
	rg.DELETE("/printers/:printerId", e.deletePrinter)
	rg.GET("/users", e.getUsers)
	rg.GET("/roles", e.getRoles)
	rg.GET("/shops/:shopId/users", e.getShopUsers)
	rg.POST("/shops/:shopId/users/:userId/roles", e.addUserRoleForShop)
	rg.DELETE("/shops/:shopId/users/:userId/roles/:roleId", e.deleteUserRoleFromShop)
	return rg
}

func getFromRequestCtx[T any](ctx *gin.Context, key string) (T, error) {
	x, exists := ctx.Get(key)
	if !exists {
		return *new(T), appError.NewValidationError(fmt.Errorf("request context does not contain key %s", key))
	}
	v, ok := x.(T)
	fmt.Sprintln(x)
	if !ok {
		return *new(T), appError.NewValidationError(fmt.Errorf("invalid type in request context for key %s", key))
	}
	return v, nil
}

func NewApiEnv(logic logic.ILogic, authenticator iface.IAuthenticator[*domain.OidcUser]) *ApiEnv {
	return &ApiEnv{
		logic:         logic,
		authenticator: authenticator,
	}
}
