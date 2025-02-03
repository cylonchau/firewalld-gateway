package v2

import (
	"github.com/gin-gonic/gin"

	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"
)

type MasqueradeRouter struct{}

func (this *MasqueradeRouter) RegisterPortAPI(g *gin.RouterGroup) {
	masqueradeGroup := g.Group("/masquerade")

	masqueradeGroup.PUT("/", this.enableInPermanent)
	masqueradeGroup.DELETE("/", this.disableInPermanent)
	masqueradeGroup.GET("/", this.queryInPermanent)
}

// enableInPermanent godoc
// @Summary Enable masqerade at firewall permanent.
// @Description Enable masqerade at firewall permanent.
// @Tags firewalld masquerade
// @Accept  json
// @Produce json
// @Param query body query.Query  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/masquerade [put]
func (this *MasqueradeRouter) enableInPermanent(c *gin.Context) {

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

	if err := dbusClient.EnablePermanentMasquerade(query.Zone); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, nil)
}

// disableInPermanent godoc
// @Summary Disable masqerade at firewall permanent.
// @Description Disable masqerade at firewall permanent.
// @Tags firewalld masquerade
// @Accept  json
// @Produce json
// @Param query  body  query.Query  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/masquerade [delete]
func (this *MasqueradeRouter) disableInPermanent(c *gin.Context) {

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

	if err := dbusClient.DisablePermanentMasquerade(query.Zone); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, nil)
}

// queryInPermanent godoc
// @Summary Get status of masqurade at firewalld permanent.
// @Description Get status of masqurade at firewalld permanent.
// @Tags firewalld masquerade
// @Accept json
// @Produce json
// @Param ip  query  string true "body"
// @Security BearerAuth
// @Success 200 {object} []interface{}
// @Router /fw/v2/masquerade [get]
func (this *MasqueradeRouter) queryInPermanent(c *gin.Context) {

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

	isenable, err := dbusClient.QueryPermanentMasquerade(query.Zone)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	if isenable == false {
		api_query.SuccessResponse(c, api_query.NETWORK_MASQUERADE_DISABLE, isenable)
		return
	}

	api_query.SuccessResponse(c, api_query.NETWORK_MASQUERADE_ENABLE, isenable)
}
