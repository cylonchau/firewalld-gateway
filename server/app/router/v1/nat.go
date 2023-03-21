package v1

import (
	q "github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type NATRouter struct{}

func (this *NATRouter) RegisterNATRouterAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/nat")

	portGroup.POST("/", this.addForwardInRuntime)
	portGroup.GET("/", this.getForwardInRuntime)
	portGroup.DELETE("/", this.delForwardInRuntime)
}

// addForward ...
// @Summary addForward
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/add [POST]
func (this *NATRouter) addForwardInRuntime(c *gin.Context) {

	var query = &q.ForwardQuery{}
	if err := c.ShouldBind(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

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
func (this *NATRouter) getForwardInRuntime(c *gin.Context) {

	var query = &q.Query{}
	err := c.Bind(query)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	forwards, err := dbusClient.Listforwards(query.Zone)

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
func (this *NATRouter) delForwardInRuntime(c *gin.Context) {

	var query = &q.ForwardQuery{}
	if err := c.ShouldBind(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.RemoveForwardPort(query.Zone, query.Forward); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}
