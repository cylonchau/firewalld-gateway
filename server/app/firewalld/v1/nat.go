package v1

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type NATRouter struct{}

func (this *NATRouter) RegisterNATV1RouterAPI(g *gin.RouterGroup) {
	natGroup := g.Group("/nat")

	natGroup.PUT("/", this.addForwardInRuntime)
	natGroup.GET("/", this.getForwardInRuntime)
	natGroup.DELETE("/", this.delForwardInRuntime)
}

// addForwardOnRuntime godoc
// @Summary Add a nat rule at firewall runtime.
// @Description Add a nat rule at firewall runtime.
// @Tags firewalld NAT
// @Accept json
// @Produce json
// @Param query  body  query.ForwardQuery  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/nat [put]
func (this *NATRouter) addForwardInRuntime(c *gin.Context) {

	var query = &api_query.ForwardQuery{}
	if err := c.ShouldBind(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.AddForwardPort(query.Zone, query.Timeout, query.Forward); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}

// getForwardInRuntime godoc
// @Summary Get nat rules at firewalld runtime.
// @Description Get nat rules at firewalld runtime
// @Tags firewalld NAT
// @Accept  json
// @Produce json
// @Param query body query.Query  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/nat [get]
func (this *NATRouter) getForwardInRuntime(c *gin.Context) {

	var query = &api_query.Query{}
	err := c.Bind(query)

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

	forwards, err := dbusClient.Listforwards(query.Zone)

	if err != nil {
		api_query.SuccessResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, forwards)
}

// delForwardInRuntime godoc
// @Summary Remove a nat rule at firewalld runtime.
// @Description Remove a nat rule at firewalld runtime.
// @Tags firewalld NAT
// @Accept  json
// @Produce json
// @Param  query  body  query.Query  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/nat [delete]
func (this *NATRouter) delForwardInRuntime(c *gin.Context) {

	var query = &api_query.ForwardQuery{}
	if err := c.ShouldBind(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.RemoveForwardPort(query.Zone, query.Forward); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}
