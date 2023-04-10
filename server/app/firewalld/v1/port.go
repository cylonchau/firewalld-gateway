package v1

import (
	code "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

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

	var query = &code.Query{}
	err := c.Bind(query)

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

	port, err := dbusClient.GetPort(query.Zone)

	if err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	if len(port) <= 0 {
		code.NotFount(c, code.ErrPortNotFount, port)
		return
	}

	code.SuccessResponse(c, code.OK, port)
}

// AddPort ...
// @Summary AddPort
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/add [POST]
func (this *PortRouter) addInRuntime(c *gin.Context) {

	var query = &code.PortQuery{}
	if err := c.BindJSON(query); err != nil {
		code.APIResponse(c, err, nil)
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
		code.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.AddPort(&query.Port, query.Zone, query.Timeout); err != nil {
		code.APIResponse(c, err, nil)
		return
	}
	code.SuccessResponse(c, code.OK, query)
}

// AddPort ...
// @Summary AddPort
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/add [POST]
func (this *PortRouter) removeInRuntime(c *gin.Context) {

	var query = &code.PortQuery{}
	if err := c.BindJSON(query); err != nil {
		code.APIResponse(c, err, nil)
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
		code.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.RemovePort(&query.Port, query.Zone); err != nil {
		code.APIResponse(c, err, nil)
		return
	}
	code.SuccessResponse(c, code.OK, query)
}
