package v1

import (
	"github.com/gin-gonic/gin"

	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
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

func (this *DashboardRouter) getRuntimeStatus(c *gin.Context) {

	var query = &api_query.Query{}
	err := c.BindQuery(query)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
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
			if services, err := dbusClient.ListServices(); err == nil {
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
		api_query.SuccessResponse(c, api_query.OK, status)
	} else {
		api_query.SuccessResponse(c, api_query.OK, err)
	}

}

func (this *DashboardRouter) getHostPie(c *gin.Context) {
	var (
		encounter    error
		hostClassify []*model.Classify
	)

	if hostClassify, encounter = model.HostClassify(); encounter == nil {
		api_query.SuccessResponse(c, api_query.OK, hostClassify)
		return
	}
	api_query.SuccessResponse(c, api_query.ErrDashboardFailed, nil)
}

func (this *DashboardRouter) getHostPanel(c *gin.Context) {
	var hostCount, tagCount, templateCount int64

	hostCount = model.HostCounter()
	tagCount = model.TagCounter()
	templateCount = model.TemplateCounter()
	status := make(map[string]interface{})
	status["hosts"] = hostCount
	status["tags"] = tagCount
	status["templates"] = templateCount
	api_query.SuccessResponse(c, api_query.OK, status)
	return
}
