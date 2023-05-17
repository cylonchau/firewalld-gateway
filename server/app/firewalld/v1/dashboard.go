package v1

import (
	"github.com/gin-gonic/gin"

	code "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

type DashboardRouter struct{}

func (this *DashboardRouter) RegisterPortAPI(g *gin.RouterGroup) {
	dashboardGroup := g.Group("/dashboard")
	dashboardGroup.GET("/", this.getRuntimeStatus)
	dashboardGroup.GET("/panel", this.getHostPanel)
	dashboardGroup.GET("/pie", this.getHostPie)

}

// getRuntimeStatus ...
// @Summary getRuntimeStatus
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/dashboard [GET]
func (this *DashboardRouter) getRuntimeStatus(c *gin.Context) {

	var query = &code.Query{}
	err := c.BindQuery(query)

	if err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		code.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()
	defaultPolicy := dbusClient.GetDefaultPolicy()
	defaultZone := dbusClient.GetDefaultZone()
	var richCount, portCount, serviceCount int
	var natStatus bool

	if richs, err := dbusClient.GetRichRules(defaultZone); err == nil {
		richCount = len(richs)
		if ports, err := dbusClient.GetPorts(defaultZone); err == nil {
			portCount = len(ports)
			if services, err := dbusClient.GetServices(); err == nil {
				serviceCount = len(services)
				if b, err := dbusClient.QueryMasquerade(defaultZone); err == nil {
					natStatus = b
				}
			}
		}
	}

	if err == nil {
		status := make(map[string]interface{})
		status["default_zone"] = defaultZone
		status["default_policy"] = defaultPolicy
		status["nat_status"] = natStatus
		status["rich"] = richCount
		status["port"] = portCount
		status["service"] = serviceCount
		code.SuccessResponse(c, code.OK, status)
	} else {
		code.SuccessResponse(c, code.OK, err)
	}

}

// getDBStatus ...
// @Summary getDBStatus
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/dashboard/template [GET]
func (this *DashboardRouter) getHostPie(c *gin.Context) {
	var (
		encounter    error
		hostClassify []*model.Classify
	)

	if hostClassify, encounter = model.HostClassify(); encounter == nil {
		code.SuccessResponse(c, code.OK, hostClassify)
		return
	}
	code.SuccessResponse(c, code.ErrDashboardFailed, nil)
}

// getHostPie ...
// @Summary getHostPie
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/dashboard/pie [GET]
func (this *DashboardRouter) getHostPanel(c *gin.Context) {
	var hostCount, tagCount, templateCount int64

	hostCount = model.HostCounter()
	tagCount = model.TagCounter()
	templateCount = model.TemplateCounter()
	status := make(map[string]interface{})
	status["hosts"] = hostCount
	status["tags"] = tagCount
	status["templates"] = templateCount
	code.SuccessResponse(c, code.OK, status)
	return
}
