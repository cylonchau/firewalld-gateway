package v2

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"
)

type ServiceRouter struct{}

func (this *ServiceRouter) RegisterPortAPI(g *gin.RouterGroup) {
	serivceGroup := g.Group("/service")
	serivceGroup.PUT("/", this.addServicesPermanent)
	serivceGroup.GET("/", this.getServicesAtPermanent)
	serivceGroup.DELETE("/", this.deleteServicesAtPermanent)
	serivceGroup.GET("/config", this.listServices)
	serivceGroup.PUT("/config", this.newServiceAtPermanent)
}

// getServicesAtRuntime godoc
// @Summary List available service type on firewalld permanent.
// @Description List all service rule on firewalld permanent.
// @Tags firewalld service
// @Accept  json
// @Produce json
// @Param query body query.Query false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/service [get]
func (this *ServiceRouter) getServicesAtPermanent(c *gin.Context) {

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

	services, err := dbusClient.GetPermanentServices()

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

// deleteServicesAtPermanent godoc
// @Summary Remove a service rule on firewalld premanent.
// @Description Remove a service rule on firewalld premanent.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.ServiceQuery  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/service [delete]
func (this *ServiceRouter) deleteServicesAtPermanent(c *gin.Context) {

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

	if err := dbusClient.RemovePermanentService(serviceQuery.Service); err != nil {
		query.APIResponse(c, err, nil)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// addServicesPermanent godoc
// @Summary Permanently add service to list of services used in zone.
// @Description Permanently add service to list of services used in zone.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.ServiceQuery  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/service [put]
func (this *ServiceRouter) addServicesPermanent(c *gin.Context) {

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
	err = dbusClient.AddPermanentService(serviceQuery.Zone, serviceQuery.Service)
	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// listServicesAtRuntime godoc
// @Summary Return list of service names (s) in runtime configuration.
// @Description Return list of service names (s) in runtime configuration.
// @Tags firewalld service
// @Accept  json
// @Produce json
// @Param query body query.Query false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/service/config [get]
func (this *ServiceRouter) listServices(c *gin.Context) {

	var serviceQuery = &query.Query{}
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

	services, err := dbusClient.ListServices()
	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		query.NotFount(c, errors.New("list of available services is not found."), nil)
		return
	}

	query.SuccessResponse(c, query.OK, services)
}

// newServiceAtPermanent godoc
// @SummaryAdd service with given settings into permanent configuration.
// @DescriptionAdd service with given settings into permanent configuration.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.ServiceSettingQuery  false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/service/config [put]
func (this *ServiceRouter) newServiceAtPermanent(c *gin.Context) {

	var serviceSettingQuery = &query.ServiceSettingQuery{}

	if err := c.BindJSON(serviceSettingQuery); err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(serviceSettingQuery.Host)
	if err != nil {
		query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	err = dbusClient.AddNewService(serviceSettingQuery.ServiceName, serviceSettingQuery.Setting)
	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}
	query.SuccessResponse(c, query.OK, serviceSettingQuery.Setting)
}
