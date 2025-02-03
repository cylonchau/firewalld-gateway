package v2

import (
	"fmt"

	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type PortV2Router struct{}

func (this *PortV2Router) RegisterPortV2API(g *gin.RouterGroup) {
	portGroup := g.Group("/ports")
	portGroup.GET("/", this.getOnPermanent)
	portGroup.PUT("/", this.addOnPermanent)
	portGroup.DELETE("/", this.removeOnPermanent)
}

// getOnPermanent godoc
// @Summary Get port rule on firewalld permanent.
// @Description Get port rule on firewalld permanent.
// @Tags firewalld port
// @Accept json
// @Produce json
// @Param   ip  query  string true "body"
// @Security BearerAuth
// @Success 200 {object} []interface{}
// @Router /fw/v2/ports [get]
func (this *PortV2Router) getOnPermanent(c *gin.Context) {

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

	port, err := dbusClient.PermanentGetPort(query.Zone)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	if len(port) <= 0 {
		api_query.NotFount(c, api_query.ErrPortNotFount, port)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, port)
}

// addOnPermanent godoc
// @Summary Add port rule on firewalld permanent.
// @Description Add port rule on firewalld permanent.
// @Tags firewalld port
// @Accept json
// @Produce json
// @Param query body query.Query  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/ports [put]
func (this *PortV2Router) addOnPermanent(c *gin.Context) {

	var query = &api_query.PortQuery{}
	if err := c.BindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
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
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.PermanentAddPort(fmt.Sprintf("%s/%s", query.Port.Port, query.Port.Protocol), query.Zone); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}

// removeOnPermanent godoc
// @Summary Delete a firewalld port rule on permanent.
// @Description Delete a firewalld port rule on permanent.
// @Tags firewalld port
// @Accept  json
// @Produce json
// @Param query  body  query.PortQuery   false "body"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /fw/v2/ports [delete]
func (this *PortV2Router) removeOnPermanent(c *gin.Context) {

	var query = &api_query.PortQuery{}
	if err := c.BindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
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
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.PermanentRemovePort(fmt.Sprintf("%s/%s", query.Port.Port, query.Port.Protocol), query.Zone); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}
