package v1

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

// swagger_annotations.go

type MasqueradeRouter struct{}

func (this *MasqueradeRouter) RegisterPortAPI(g *gin.RouterGroup) {
	masqueradeGroup := g.Group("/masquerade")

	masqueradeGroup.PUT("/", this.enableInRuntime)
	masqueradeGroup.DELETE("/", this.disableInRuntime)
	masqueradeGroup.GET("/", this.queryInRuntime)
}

// enableInRuntime godoc
// @Summary Enable masqerade on firewalld runtime.
// @Description Enable masqerade on firewalld runtime.
// @Tags firewalld masquerade
// @Accept  json
// @Produce json
// @Param  query  body  query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/masquerade [put]
func (this *MasqueradeRouter) enableInRuntime(c *gin.Context) {

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
	if err := dbusClient.EnableMasquerade(query.Zone, query.Timeout); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, nil)
}

// disableInRuntime godoc
// @Summary Disable masqerade on firewalld runtime.
// @Description Disable masqerade on firewalld runtime.
// @Tags firewalld masquerade
// @Accept  json
// @Produce json
// @Param query  body  query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/masquerade [delete]
func (this *MasqueradeRouter) disableInRuntime(c *gin.Context) {

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
	if err := dbusClient.DisableMasquerade(query.Zone); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, nil)
}

// queryInRuntime godoc
// @Summary Get nat rule list on firewalld runtime.
// @Description Get nat rule list on firewalld runtime.
// @Tags firewalld masquerade
// @Accept  json
// @Produce json
// @Param ip  query  string true "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} []interface{}
// @Router /fw/v1/masquerade [get]
func (this *MasqueradeRouter) queryInRuntime(c *gin.Context) {

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

	isenable, err := dbusClient.QueryMasquerade(query.Zone)

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
