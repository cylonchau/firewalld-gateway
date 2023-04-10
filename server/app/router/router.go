package router

import (
	"github.com/cylonchau/firewalld-gateway/server/app/firewalld/host"
	"github.com/cylonchau/firewalld-gateway/server/app/firewalld/tag"
	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/server/app/auth"
	"github.com/cylonchau/firewalld-gateway/server/app/auther"
	fv1 "github.com/cylonchau/firewalld-gateway/server/app/firewalld/v1"
	fv2 "github.com/cylonchau/firewalld-gateway/server/app/firewalld/v2"
	fv3 "github.com/cylonchau/firewalld-gateway/server/app/firewalld/v3"
)

func RegisteredRouter(e *gin.Engine) {
	e.Handle("GET", "ping", ping)

	firewallAPI := e.Group("/fw")
	authAPI := e.Group("/auth")
	firewallAPI.Use(auther.JWTAuthMiddleware())

	fv1Group := firewallAPI.Group("/v1")
	fv2Group := firewallAPI.Group("/v2")
	fv3Group := firewallAPI.Group("/v3")

	tagGroup := firewallAPI.Group("/tag")
	hostGroup := firewallAPI.Group("/host")

	portRouter := &fv1.PortRouter{}
	portRouter.RegisterPortAPI(fv1Group)

	masqueradeRouter := &fv1.MasqueradeRouter{}
	masqueradeRouter.RegisterPortAPI(fv1Group)

	natRouter := &fv1.NATRouter{}
	natRouter.RegisterNATRouterAPI(fv1Group)

	richRuleRouter := &fv1.RichRuleRouter{}
	richRuleRouter.RegisterPortAPI(fv1Group)

	serviceRouter := &fv1.ServiceRouter{}
	serviceRouter.RegisterPortAPI(fv1Group)

	natv2Router := &fv2.NatRouter{}
	natv2Router.RegisterPortAPI(fv2Group)

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

	// auth route
	authRouter := &auth.Auth{}
	authRouter.RegisterUserAPI(authAPI)

	if !config.CONFIG.MySQL.IsEmpty() || !config.CONFIG.SQLite.IsEmpty() {
		tagRouter := &tag.Tag{}
		tagRouter.RegisterTagAPI(tagGroup)

		hostRouter := &host.Host{}
		hostRouter.RegisterHostAPI(hostGroup)

		asyncHostRouter := &host.AsyncHost{}
		asyncHostRouter.RegisterAsyncHostAPI(hostGroup)
	}
}
