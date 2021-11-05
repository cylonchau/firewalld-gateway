package router

import (
	v1 "firewall-api/cmd/app/router/v1"
	v2 "firewall-api/cmd/app/router/v2"
	"firewall-api/config"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(e *gin.Engine) {
	firewall_api := e.Group("/fw")
	v1Group := firewall_api.Group("/v1")
	v2Group := firewall_api.Group("/v2")

	portRouter := &v1.PortRouter{}
	portRouter.RegisterPortAPI(v1Group)

	masqueradeRouter := &v1.MasqueradeRouter{}
	masqueradeRouter.RegisterPortAPI(v1Group)

	natRouter := &v1.NatRouter{}
	natRouter.RegisterPortAPI(v1Group)

	richRuleRouter := &v1.RichRuleRouter{}
	richRuleRouter.RegisterPortAPI(v1Group)

	serviceRouter := &v1.ServiceRouter{}
	serviceRouter.RegisterPortAPI(v1Group)

	natv2Router := &v2.NatRouter{}
	natv2Router.RegisterPortAPI(v2Group)

	settingRouter := v2.SettingRouter{}
	settingRouter.RegisterPortAPI(v2Group)

	if !config.CONFIG.Mysql.IsEmpty() {
		storageRouter := v1.StoageRouter{}
		storageRouter.RegisterStoageAPI(v1Group)
	}
}
