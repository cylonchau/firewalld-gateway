package v1

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type ServiceRouter struct{}

func (this *ServiceRouter) RegisterPortAPI(g *gin.RouterGroup) {
	serivceGroup := g.Group("/service")
	serivceGroup.GET("/", this.getServicesAtRuntime)
	serivceGroup.DELETE("/", this.deleteServicesAtRuntime)
	serivceGroup.PUT("/", this.addServicesAtRuntime)

}

// getServicesAtRuntime godoc
// @Summary List available service type on firewalld.
// @Description List all service rule on firewalld.
// @Tags firewalld service
// @Accept  json
// @Produce json
// @Param query body query.Query false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/service [get]
func (this *ServiceRouter) getServicesAtRuntime(c *gin.Context) {

	var rich = &api_query.Query{}
	err := c.BindQuery(rich)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(rich.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	services, err := dbusClient.GetServices()

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		api_query.NotFount(c, api_query.ErrServiceNotFount, services)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, services)
}

// addServicesAtRuntime godoc
// @Summary Add a service rule on firewalld runtime.
// @Description Add a service rule on firewalld runtime.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/service [put]
func (this *ServiceRouter) addServicesAtRuntime(c *gin.Context) {

	var query = &api_query.Query{}
	err := c.BindJSON(query)

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
	err = dbusClient.AddService(query.Zone, query.Service, query.Timeout)
	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, nil)
}

// deleteServicesAtRuntime godoc
// @Summary Remove a service rule on firewalld runtime.
// @Description Remove a service rule on firewalld runtime.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/service [delete]
func (this *ServiceRouter) deleteServicesAtRuntime(c *gin.Context) {

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

	if err := dbusClient.RemoveService(query.Zone, query.Service); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, nil)
}
