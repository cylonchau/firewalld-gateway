package v2

import (
	q "github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type NatRouter struct{}

func (this *NatRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/nat")

	portGroup.POST("/", this.addForwardAtPermanent)
	portGroup.GET("/", this.getForwardAtPermanent)
	portGroup.DELETE("/", this.delForwardAtPermanent)
}

// addForwardAtPermanent ...
// @Summary addForwardAtPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/port/add [POST]
func (this *NatRouter) addForwardAtPermanent(c *gin.Context) {

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

	if err = dbusClient.AddPermanentForwardPort(query.Zone, query.Forward); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}

// getForwardAtPermanent ...
// @Summary getForwardAtPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/get [GET]
func (this *NatRouter) getForwardAtPermanent(c *gin.Context) {

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

	forwards, err := dbusClient.PermanentGetForwardPort(query.Zone)

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

// delForwardAtPermanent ...
// @Summary delForwardAtPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/port/delete [DELETE]
func (this *NatRouter) delForwardAtPermanent(c *gin.Context) {

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

	if err = dbusClient.RemovePermanentForwardPort(query.Zone, query.Forward); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query)
}
