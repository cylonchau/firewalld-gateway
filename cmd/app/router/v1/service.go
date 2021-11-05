package v1

import (
	"errors"
	"firewall-api/code"
	"firewall-api/utils/dbus"
	q "firewall-api/utils/query"

	"github.com/gin-gonic/gin"
)

type ServiceRouter struct{}

func (this *ServiceRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/service")
	portGroup.GET("/get", this.getServicesAtRuntime)
	portGroup.POST("/add", this.addServicesAtRuntime)
	portGroup.POST("/new", this.newServiceAtPermanent)
	portGroup.GET("/list", this.listServicesAtRuntime)
	portGroup.DELETE("/delete", this.deleteServicesAtRuntime)

}

// getServicesAtRuntime ...
// @Summary getServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/get [GET]
func (this *ServiceRouter) getServicesAtRuntime(c *gin.Context) {

	var rich = &q.Query{}
	err := c.BindQuery(rich)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(rich.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	services, err := dbusClient.GetService(rich.Zone)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		q.NotFount(c, code.ErrServiceNotFount, services)
		return
	}

	q.SuccessResponse(c, code.OK, services)
}

// addServicesAtRuntime ...
// @Summary addServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/add [GET]
func (this *ServiceRouter) addServicesAtRuntime(c *gin.Context) {

	var query = &q.Query{}
	err := c.BindJSON(query)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	_, err = dbusClient.AddService(query.Zone, query.Service, query.Timeout)
	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, nil)
}

// listServicesAtRuntime ...
// @Summary listServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/list [GET]
func (this *ServiceRouter) listServicesAtRuntime(c *gin.Context) {

	var query = &q.Query{}
	err := c.Bind(query)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	services, err := dbusClient.ListServices()
	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		q.NotFount(c, errors.New("list of available services is not found."), nil)
		return
	}

	q.SuccessResponse(c, code.OK, services)
}

// deleteServicesAtRuntime ...
// @Summary deleteServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/list [DELETE]
func (this *ServiceRouter) deleteServicesAtRuntime(c *gin.Context) {

	var query = &q.Query{}
	err := c.Bind(query)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if err := dbusClient.RemoveService(query.Zone, query.Service); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, nil)
}

// newServiceAtPermanent ...
// @Summary newServiceAtPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/list [POST]
func (this *ServiceRouter) newServiceAtPermanent(c *gin.Context) {

	var query = &q.ServiceSettingQuery{}

	if err := c.BindJSON(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	err = dbusClient.NewService(query.Name, query.Setting)
	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query.Setting)
}
