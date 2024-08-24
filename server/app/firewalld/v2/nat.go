package v2

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type NATRouter struct{}

func (this *NATRouter) RegisterNATV2API(g *gin.RouterGroup) {
	portGroup := g.Group("/nat")

	portGroup.PUT("/", this.addForwardInPermanent)
	portGroup.GET("/", this.getForwardInPermanent)
	portGroup.DELETE("/", this.delForwardInPermanent)
}

// addForwardInPermanent godoc
// @Summary Add a nat rule at firewall permanent.
// @Description Add a nat rule at firewall permanent.
// @Tags firewalld NAT
// @Accept json
// @Produce json
// @Param query  body  query.ForwardQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/nat [put]
func (this *NATRouter) addForwardInPermanent(c *gin.Context) {

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

	if err = dbusClient.AddPermanentForwardPort(query.Zone, query.Forward); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}

// getForwardInPermanent godoc
// @Summary Get nat rules at firewalld permanent.
// @Description Get nat rules at firewalld permanent.
// @Tags firewalld NAT
// @Accept json
// @Produce json
// @Param ip query string true "ip"
// @Param zone query string false "zone"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/nat [get]
func (this *NATRouter) getForwardInPermanent(c *gin.Context) {

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

	forwards, err := dbusClient.PermanentGetForwardPort(query.Zone)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, forwards)
}

// delForwardInPermanent godoc
// @Summary Remove a nat rule at firewalld permanent.
// @Description Remove a nat rule at firewalld permanent.
// @Tags firewalld NAT
// @Accept json
// @Produce json
// @Param query  body  query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/nat [delete]
func (this *NATRouter) delForwardInPermanent(c *gin.Context) {

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

	if err = dbusClient.RemovePermanentForwardPort(query.Zone, query.Forward); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}
