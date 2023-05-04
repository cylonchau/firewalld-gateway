package v1

import (
	code "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type MasqueradeRouter struct{}

func (this *MasqueradeRouter) RegisterPortAPI(g *gin.RouterGroup) {
	masqueradeGroup := g.Group("/masquerade")

	masqueradeGroup.PUT("/", this.enableInRuntime)
	masqueradeGroup.DELETE("/", this.disableInRuntime)
	masqueradeGroup.GET("/", this.queryInRuntime)
	masqueradeGroup.PUT("/permanent", this.enableInPermanent)
	masqueradeGroup.DELETE("/permanent", this.disableInPermanent)
	masqueradeGroup.GET("/query", this.queryInPermanent)
}

// enableInRuntime ...
// @Summary enableInRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/enable [GET]
func (this *MasqueradeRouter) enableInRuntime(c *gin.Context) {

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
	if err := dbusClient.EnableMasquerade(query.Zone, query.Timeout); err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	code.SuccessResponse(c, code.OK, nil)
}

// disableInRuntime ...
// @Summary disableInRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/disable [GET]
func (this *MasqueradeRouter) disableInRuntime(c *gin.Context) {

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
	if err := dbusClient.DisableMasquerade(query.Zone); err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	code.SuccessResponse(c, code.OK, nil)
}

// queryInRuntime ...
// @Summary queryInRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/query [GET]
func (this *MasqueradeRouter) queryInRuntime(c *gin.Context) {

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

	isenable, err := dbusClient.QueryMasquerade(query.Zone)

	if err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	if isenable == false {
		code.SuccessResponse(c, code.NETWORK_MASQUERADE_DISABLE, isenable)
		return
	}

	code.SuccessResponse(c, code.NETWORK_MASQUERADE_ENABLE, isenable)
}

// enableInPermanent ...
// @Summary enableInPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/enableinpermanent [GET]
func (this *MasqueradeRouter) enableInPermanent(c *gin.Context) {

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

	if err := dbusClient.EnablePermanentMasquerade(query.Zone); err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	code.SuccessResponse(c, code.OK, nil)
}

// disableInPermanent ...
// @Summary disableInPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/disablepermanent [GET]
func (this *MasqueradeRouter) disableInPermanent(c *gin.Context) {

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

	if err := dbusClient.DisablePermanentMasquerade(query.Zone); err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	code.SuccessResponse(c, code.OK, nil)
}

// queryInPermanent ...
// @Summary queryInPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/querypermanent [GET]
func (this *MasqueradeRouter) queryInPermanent(c *gin.Context) {

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

	isenable, err := dbusClient.QueryPermanentMasquerade(query.Zone)

	if err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	if isenable == false {
		code.SuccessResponse(c, code.NETWORK_MASQUERADE_DISABLE, isenable)
		return
	}

	code.SuccessResponse(c, code.NETWORK_MASQUERADE_ENABLE, isenable)
}
