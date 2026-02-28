package api

import (
	"path"
	"reflect"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/uoul/go-auth"
	"github.com/uoul/klcs/backend/core/api/dto"
	"github.com/uoul/klcs/backend/core/apperror"
	"github.com/uoul/klcs/backend/core/config"
	"github.com/uoul/klcs/backend/core/domain"
	"github.com/uoul/klcs/backend/core/logic"
	"github.com/uoul/klcs/backend/core/services"
)

type Api struct {
	logic         *logic.Logic
	authenticator auth.IAuthenticator[domain.OidcUser]
	printService  *services.PrintService
	version       string

	mode       string
	wwwRootDir string
	cors       config.CorsConfig
	oidc       config.OidcConfig
	uiConfig   config.UiConfig
}

func (e *Api) Run(api string) error {
	// Get new router
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		cors.New(cors.Config{
			AllowWildcard:   true,
			AllowWebSockets: true,
			AllowOrigins:    e.cors.Origins,
			AllowMethods:    e.cors.Methods,
			AllowHeaders:    e.cors.Headers,
		}),
		static.Serve("/", static.LocalFile(e.wwwRootDir, true)),
	)
	// Handle if no route match --> redirect to index.html
	router.NoRoute(func(c *gin.Context) {
		c.File(path.Join(e.wwwRootDir, "index.html"))
	})
	// Setup public api
	publicRg := router.Group("/public")
	publicRg.Use(
		e.errorTranslation(),
	)
	publicRg.GET("/settings", e.getUiSettings)
	apiV1 := router.Group("/api/v1")
	// Setup global middleware
	apiV1.Use(
		e.errorTranslation(),
		e.checkUserLoggedIn(),
		e.updateUserByOidcData(),
	)
	// Setup routergroups
	e.setupSysAdminRg(apiV1, "/admin")
	e.setupUserRg(apiV1, "")
	e.setupAccountManagerRg(apiV1, "/accounts")
	e.setupPrinterApi(&router.RouterGroup, "/api/v1/printers")
	// Run api
	return router.Run(api)
}

func (e *Api) setupSysAdminRg(router *gin.RouterGroup, prefix string) *gin.RouterGroup {
	rg := router.Group(prefix)
	rg.Use(
		e.checkOidcRole(e.oidc.Roles.SysAdmin),
	)
	// assign handlers
	rg.POST("/shops", e.createShop)
	rg.GET("/shops", e.getShops)
	rg.PATCH("/shops/:shopId", e.updateShop)
	rg.DELETE("/shops/:shopId", e.deleteShop)
	return rg
}

func (e *Api) setupAccountManagerRg(router *gin.RouterGroup, prefix string) {
	rg := router.Group(prefix)
	rg.Use(
		e.checkOidcRole(e.oidc.Roles.AccountManager),
	)
	rg.GET("", e.getAccounts)
	rg.POST("", e.createAccount)
	rg.PATCH("/:accountId", e.updateAccount)
	rg.DELETE("/:accountId/balance", e.closeAccount)
	rg.POST("/:accountId/balance", e.postToAccount)
}

func (a *Api) setupPrinterApi(router *gin.RouterGroup, prefix string) {
	rg := router.Group(prefix)
	rg.GET("/:printerId/jobs", a.getPrintJobs)
	rg.POST("/:printerId/jobs/acknowledgement/:transactionId", a.acknowledgePrintJob)
}

func (e *Api) setupUserRg(router *gin.RouterGroup, prefix string) *gin.RouterGroup {
	rg := router.Group(prefix)
	// seller api
	rg.GET("/shops", e.getShopsForUser)
	rg.GET("/shops/:shopId", e.getShopdetailsForUser)
	rg.POST("/orders", e.checkout)
	rg.GET("/accounts/:accountId", e.getAccountDetails)
	rg.GET("/history", e.getHistoryForUser)
	rg.POST("/orders/:transactionId/printjob", e.reprint)
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

func WithApiWwwRootDir(dir string) func(*Api) {
	return func(a *Api) {
		a.wwwRootDir = dir

	}
}

func WithCorsOrigins(v []string) func(*Api) {
	return func(a *Api) {
		a.cors.Origins = v
	}
}

func WithCorsMethods(v []string) func(*Api) {
	return func(a *Api) {
		a.cors.Methods = v
	}
}

func WithCorsHeaders(v []string) func(*Api) {
	return func(a *Api) {
		a.cors.Headers = v
	}
}

func WithReleaseMode() func(*Api) {
	return func(a *Api) {
		a.mode = gin.ReleaseMode
	}
}

func NewApi(version string, corsCfg config.CorsConfig, oidcCfg config.OidcConfig, uiConfig config.UiConfig, logic *logic.Logic, authenticator auth.IAuthenticator[domain.OidcUser], printService *services.PrintService, opts ...func(*Api)) *Api {
	// Create Instance
	a := &Api{
		logic:         logic,
		authenticator: authenticator,
		printService:  printService,
		cors:          corsCfg,
		oidc:          oidcCfg,
		version:       version,
		uiConfig:      uiConfig,

		wwwRootDir: "wwwroot",
		mode:       gin.DebugMode,
	}
	// Apply Options
	for _, o := range opts {
		o(a)
	}
	// Return Api
	return a
}

func (a *Api) NewErrorResponse(err apperror.IAppError) *dto.ErrorResponse {
	var msg string
	if a.mode == gin.DebugMode {
		msg = err.Error()
	}
	t := reflect.TypeOf(err)
	// Handle if Error is a pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return &dto.ErrorResponse{
		Code:         t.Name(),
		DebugMessage: msg,
	}
}
