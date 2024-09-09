package v1

import (
	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
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

// addServicesAtRuntime godoc
// @Summary Add a service rule on firewalld runtime.
// @Description Add a service rule on firewalld runtime.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.ServiceQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/service [put]
func (this *ServiceRouter) addServicesAtRuntime(c *gin.Context) {

	var serviceQuery = &query.ServiceQuery{}
	err := c.BindJSON(serviceQuery)

	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(serviceQuery.Ip)
	if err != nil {
		query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()
	err = dbusClient.AddServiceRuntime(serviceQuery.Zone, serviceQuery.Service, serviceQuery.Timeout)
	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// deleteServicesAtRuntime godoc
// @Summary Remove a service rule on firewalld runtime.
// @Description Remove a service rule on firewalld runtime.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.ServiceQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/service [delete]
func (this *ServiceRouter) deleteServicesAtRuntime(c *gin.Context) {

	var serviceQuery = &query.ServiceQuery{}
	err := c.Bind(serviceQuery)

	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(serviceQuery.Ip)
	if err != nil {
		query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err := dbusClient.RemoveRuntimeService(serviceQuery.Zone, serviceQuery.Service); err != nil {
		query.APIResponse(c, err, nil)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// getServicesAtRuntime godoc
// @Summary List available service type on firewalld runtime.
// @Description List all service rule on firewalld runtime.
// @Tags firewalld service
// @Accept  json
// @Produce json
// @Param query body query.ServiceQuery false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/service [get]
func (this *ServiceRouter) getServicesAtRuntime(c *gin.Context) {

	var serviceQuery = &query.Query{}
	err := c.BindQuery(serviceQuery)
	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(serviceQuery.Ip)
	if err != nil {
		query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	services, err := dbusClient.GetRuntimeServices(serviceQuery.Zone)

	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		query.NotFount(c, query.ErrServiceNotFount, services)
		return
	}

	query.SuccessResponse(c, query.OK, services)
}
