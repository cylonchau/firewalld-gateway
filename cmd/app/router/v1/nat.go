package v1

import (
	"firewall-api/code"
	"firewall-api/utils/dbus"
	q "firewall-api/utils/query"
	"github.com/gin-gonic/gin"
)

type NatRouter struct{}

func (this *NatRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/nat")

	portGroup.POST("/add", this.addForwardInRuntime)
	portGroup.GET("/get", this.getForwardInRuntime)
	portGroup.DELETE("/delete", this.delForwardInRuntime)
}

// addForward ...
// @Summary addForward
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/add [POST]
func (this *NatRouter) addForwardInRuntime(c *gin.Context) {

	var query = &q.ForwardQuery{}
	if err := c.ShouldBind(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if err = dbusClient.AddForwardPort(query.Zone, query.Timeout, query.Forward); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}

// GetPort ...
// @Summary GetPort
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/get [GET]
func (this *NatRouter) getForwardInRuntime(c *gin.Context) {

	var query = &q.Query{}
	err := c.Bind(query)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	forwards, err := dbusClient.GetForwardPort(query.Zone)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	if len(forwards) <= 0 {
		q.NotFount(c, code.ErrForwardNotFount, nil)
		return
	}
	q.SuccessResponse(c, code.OK, forwards)
}

// delForwardInRuntime ...
// @Summary delForwardInRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/add [DELETE]
func (this *NatRouter) delForwardInRuntime(c *gin.Context) {

	var query = &q.ForwardQuery{}
	if err := c.ShouldBind(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if err = dbusClient.RemoveForwardPort(query.Zone, query.Forward); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}
