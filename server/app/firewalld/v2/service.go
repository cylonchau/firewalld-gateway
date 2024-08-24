package v2

import (
	"errors"

	"github.com/gin-gonic/gin"

	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"
)

type ServiceRouter struct{}

func (this *ServiceRouter) RegisterPortAPI(g *gin.RouterGroup) {
	serivceGroup := g.Group("/service")
	serivceGroup.PUT("/", this.newServiceAtPermanent)
	serivceGroup.GET("/types", this.listServicesAtRuntime)
}

// addServicesAtRuntime godoc
// @Summary Add a service rule at firewalld permanent.
// @Description Add a service rule at firewalld permanent.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.ServiceSettingQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/service [put]
func (this *ServiceRouter) newServiceAtPermanent(c *gin.Context) {

	var query = &api_query.ServiceSettingQuery{}

	if err := c.BindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Host)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	err = dbusClient.AddNewService(query.ServiceName, query.Setting)
	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query.Setting)
}

// listServicesAtRuntime godoc
// @Summary List available service type on firewalld.
// @Description List all service rule on firewalld.
// @Tags firewalld service
// @Accept  json
// @Produce json
// @Param query body query.Query false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/service [get]
func (this *ServiceRouter) listServicesAtRuntime(c *gin.Context) {

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

	services, err := dbusClient.GetServices()
	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		api_query.NotFount(c, errors.New("list of available services is not found."), nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, services)
}
