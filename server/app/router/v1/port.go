package v1

import (
	q "github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type PortRouter struct{}

func (this *PortRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/port")
	portGroup.GET("/", this.getInRuntime)
	portGroup.POST("/", this.addInRuntime)
	portGroup.DELETE("/", this.removeInRuntime)
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

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

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

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.AddPort(&query.Port, query.Zone, query.Timeout); err != nil {
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

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.RemovePort(&query.Port, query.Zone); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}
