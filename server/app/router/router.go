package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/server/app/audit"
	"github.com/cylonchau/firewalld-gateway/server/app/auth"
	"github.com/cylonchau/firewalld-gateway/server/app/firewalld/host"
	"github.com/cylonchau/firewalld-gateway/server/app/firewalld/tag"
	"github.com/cylonchau/firewalld-gateway/server/app/firewalld/template"
	"github.com/cylonchau/firewalld-gateway/server/app/firewalld/token"
	fv1 "github.com/cylonchau/firewalld-gateway/server/app/firewalld/v1"
	fv2 "github.com/cylonchau/firewalld-gateway/server/app/firewalld/v2"
	fv3 "github.com/cylonchau/firewalld-gateway/server/app/firewalld/v3"
	"github.com/cylonchau/firewalld-gateway/server/app/middlewares"
	"github.com/cylonchau/firewalld-gateway/server/app/sso"
	user "github.com/cylonchau/firewalld-gateway/server/app/users"

	_ "github.com/cylonchau/firewalld-gateway/docs"
)

func RegisteredRouter(e *gin.Engine) {

	e.Handle("GET", "/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	e.Handle("GET", "ping", ping)

	ssoGroup := e.Group("/sso")
	securityAPIGroup := e.Group("/security")
	firewallAPIGroup := e.Group("/fw")
	firewallAPIGroup.Use(middlewares.JWTAuthMiddleware())
	securityAPIGroup.Use(middlewares.JWTAuthMiddleware())

	tagGroup := firewallAPIGroup.Group("/tag")
	hostGroup := firewallAPIGroup.Group("/host")
	templateGroup := firewallAPIGroup.Group("/template")

	/* Authentication  */
	// Users
	userAPI := securityAPIGroup.Group("/users")
	// tokens
	tokenAPI := securityAPIGroup.Group("/tokens")
	// Audit
	auditAPI := securityAPIGroup.Group("/audit")
	// auth
	authAPI := securityAPIGroup.Group("/auth")

	/* firewall  */
	fv1Group := firewallAPIGroup.Group("/v1")
	fv2Group := firewallAPIGroup.Group("/v2")
	fv3Group := firewallAPIGroup.Group("/v3")

	portV1Router := &fv1.PortV1Router{}
	portV2Router := &fv2.PortV2Router{}
	portV1Router.RegisterPortV1API(fv1Group)
	portV2Router.RegisterPortV2API(fv2Group)

	masqueradeRouterV1 := &fv1.MasqueradeRouter{}
	masqueradeRouterV2 := &fv2.MasqueradeRouter{}
	masqueradeRouterV1.RegisterPortAPI(fv1Group)
	masqueradeRouterV2.RegisterPortAPI(fv2Group)

	natRouterV1 := &fv1.NATRouter{}
	natRouterV1.RegisterNATV1RouterAPI(fv1Group)
	natRouterV2 := &fv2.NATRouter{}
	natRouterV2.RegisterNATV2API(fv2Group)

	richRuleRouterV1 := &fv1.RichRuleV1Router{}
	richRuleRouterV2 := &fv2.RichRuleV2Router{}
	richRuleRouterV1.RegisterPortAPI(fv1Group)
	richRuleRouterV2.RegisterPortAPI(fv2Group)

	serviceRouterV1 := &fv1.ServiceRouter{}
	serviceRouterV2 := &fv2.ServiceRouter{}
	serviceRouterV1.RegisterPortAPI(fv1Group)
	serviceRouterV2.RegisterPortAPI(fv2Group)

	dashboardRouter := &fv1.DashboardRouter{}
	dashboardRouter.RegisterPortAPI(fv1Group)

	settingRouter := fv2.SettingRouter{}
	settingRouter.RegisterPortAPI(fv2Group)

	if config.CONFIG.AsyncProcess {
		batchPortRouter := fv3.PortRouter{}
		batchPortRouter.RegisterBatchAPI(fv3Group)

		batchSettingRouter := fv3.SettingRouter{}
		batchSettingRouter.RegisterBatchAPI(fv3Group)

		batchNATRouter := fv3.MasqueradeRouter{}
		batchNATRouter.RegisterBatchAPI(fv3Group)

		batchServiceRouter := fv3.ServiceRouter{}
		batchServiceRouter.RegisterBatchAPI(fv3Group)
	}

	if !config.CONFIG.MySQL.IsEmpty() || !config.CONFIG.SQLite.IsEmpty() {
		ssoRouter := &sso.SSO{}
		ssoRouter.RegisterUserAPI(ssoGroup)
		// auth route
		authRouter := &auth.Auth{}
		authRouter.RegisterUserAPI(authAPI)

		usersRouter := &user.User{}
		usersRouter.RegisterUserAPI(userAPI)

		tokenRouter := &Token.Token{}
		tokenRouter.RegisterTokenAPI(tokenAPI)

		tagRouter := &tag.Tag{}
		tagRouter.RegisterTagAPI(tagGroup)

		hostRouter := &host.Host{}
		hostRouter.RegisterHostAPI(hostGroup)

		asyncHostRouter := &host.AsyncHost{}
		asyncHostRouter.RegisterAsyncHostAPI(hostGroup)

		templateRouter := &template.Template{}
		templateRouter.RegisterTemplateAPI(templateGroup)

		auditRouter := &audit.Audit{}
		auditRouter.RegisterAuditAPI(auditAPI)
	}
}
