package router

import (
	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalldGateway/config"
	v1 "github.com/cylonchau/firewalldGateway/server/app/router/v1"
	v2 "github.com/cylonchau/firewalldGateway/server/app/router/v2"
	v3 "github.com/cylonchau/firewalldGateway/server/app/router/v3"
)

func RegisteredRouter(e *gin.Engine) {
	e.Handle("GET", "ping", ping)
	firewall_api := e.Group("/fw")
	v1Group := firewall_api.Group("/v1")
	v2Group := firewall_api.Group("/v2")
	v3Group := firewall_api.Group("/v3")

	portRouter := &v1.PortRouter{}
	portRouter.RegisterPortAPI(v1Group)

	masqueradeRouter := &v1.MasqueradeRouter{}
	masqueradeRouter.RegisterPortAPI(v1Group)

	natRouter := &v1.NATRouter{}
	natRouter.RegisterNATRouterAPI(v1Group)

	richRuleRouter := &v1.RichRuleRouter{}
	richRuleRouter.RegisterPortAPI(v1Group)

	serviceRouter := &v1.ServiceRouter{}
	serviceRouter.RegisterPortAPI(v1Group)

	natv2Router := &v2.NatRouter{}
	natv2Router.RegisterPortAPI(v2Group)

	settingRouter := v2.SettingRouter{}
	settingRouter.RegisterPortAPI(v2Group)

	if config.CONFIG.Async_Process {
		batchPortRouter := v3.PortRouter{}
		batchPortRouter.RegisterBatchAPI(v3Group)

		batchSettingRouter := v3.SettingRouter{}
		batchSettingRouter.RegisterBatchAPI(v3Group)

		batchNATRouter := v3.MasqueradeRouter{}
		batchNATRouter.RegisterBatchAPI(v3Group)

		batchServiceRouter := v3.ServiceRouter{}
		batchServiceRouter.RegisterBatchAPI(v3Group)
	}

	if !config.CONFIG.Mysql.IsEmpty() {
		storageRouter := v1.StoageRouter{}
		storageRouter.RegisterStoageAPI(v1Group)
	}
}
