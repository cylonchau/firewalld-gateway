package v1

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type PortV1Router struct{}

func (this *PortV1Router) RegisterPortV1API(g *gin.RouterGroup) {
	portGroup := g.Group("/ports")
	portGroup.GET("/", this.getInRuntime)
	portGroup.PUT("/", this.addInRuntime)
	portGroup.DELETE("/", this.removeInRuntime)
}

// getInRuntime godoc
// @Summary Get port rule at firewalld runtime.
// @Description Get port rule at firewalld runtime.
// @Tags firewalld port
// @Accept json
// @Produce json
// @Param   ip  query  string true "body"
// @Security BearerAuth
// @Success 200 {object} []interface{}
// @Router /fw/v1/ports [get]
func (this *PortV1Router) getInRuntime(c *gin.Context) {

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

	port, err := dbusClient.GetPorts(query.Zone)

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

// addInRuntime godoc
// @Summary Add port rule at firewalld runtime.
// @Description Add port rule at firewalld runtime.
// @Tags firewalld port
// @Accept json
// @Produce json
// @Param query  body  query.Query  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/ports [put]
func (this *PortV1Router) addInRuntime(c *gin.Context) {

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

	if err = dbusClient.AddPort(&query.Port, query.Zone, query.Timeout); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}

// removeInRuntime godoc
// @Summary Delete a firewalld port rule in runtime.
// @Description Delete a firewalld port rule in runtime.
// @Tags firewalld port
// @Accept  json
// @Produce json
// @Param query  body  query.PortQuery   false "body"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /fw/v1/ports [delete]
func (this *PortV1Router) removeInRuntime(c *gin.Context) {

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

	if err = dbusClient.RemovePort(&query.Port, query.Zone); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query)
}
