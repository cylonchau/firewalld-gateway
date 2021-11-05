package v1

import (
	"firewall-api/code"
	"firewall-api/utils/dbus"
	q "firewall-api/utils/query"
	"github.com/gin-gonic/gin"
)

type PortRouter struct{}

func (this *PortRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/port")
	portGroup.GET("/get", this.getInRuntime)
	portGroup.POST("/add", this.addInRuntime)
	portGroup.DELETE("/delete", this.removeInRuntime)
}

// GetPort ...
// @Summary GetPort
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/get [GET]
func (this *PortRouter) getInRuntime(c *gin.Context) {

	var query = &q.Query{}
	err := c.Bind(query)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	port, err := dbusClient.GetPort(query.Zone)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	if len(port) <= 0 {
		q.NotFount(c, code.ErrPortNotFount, port)
		return
	}

	q.SuccessResponse(c, code.OK, port)
}

// AddPort ...
// @Summary AddPort
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/add [POST]
func (this *PortRouter) addInRuntime(c *gin.Context) {

	var query = &q.PortQuery{}
	if err := c.BindJSON(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}
	if query.Port.Protocol == "" {
		query.Port.Protocol = "tcp"
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if err = dbusClient.AddPort(query.Port, query.Zone, query.Timeout); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}

// AddPort ...
// @Summary AddPort
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/add [POST]
func (this *PortRouter) removeInRuntime(c *gin.Context) {

	var query = &q.PortQuery{}
	if err := c.BindJSON(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}
	if query.Port.Protocol == "" {
		query.Port.Protocol = "tcp"
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if _, err = dbusClient.RemovePort(query.Port, query.Zone); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}
