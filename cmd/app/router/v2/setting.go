package v2

import (
	"firewall-api/code"
	"firewall-api/utils/dbus"
	q "firewall-api/utils/query"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ulule/deepcopier"
)

type SettingRouter struct{}

func (this *SettingRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/setting")

	portGroup.PUT("/reload", this.reload)
	portGroup.POST("/addsetting", this.addZoneSetting)
	portGroup.DELETE("/remove", this.removeZone)
	portGroup.GET("/list", this.listZone)
	portGroup.GET("/default", this.defaultZone)
	portGroup.PUT("/flush", this.flush)
	portGroup.POST("/setdefaultzone", this.setDefaultZone)

}

// reload ...
// @Summary reload
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/port/reload [PUT]
func (this *SettingRouter) reload(c *gin.Context) {

	var query = &q.Query{}
	if err := c.ShouldBind(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if err = dbusClient.Reload(); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, nil)
}

// addZoneSetting ...
// @Summary addZoneSetting
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/port/reload [PUT]
func (this *SettingRouter) addZoneSetting(c *gin.Context) {

	var query = &q.ZoneSettingQuery{}
	if err := c.BindJSON(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	var setting = &dbus.Settings{}
	var richs []string
	for _, n := range query.Setting.Rule {
		richs = append(richs, n.ToString())
	}

	deepcopier.Copy(query.Setting).To(setting)
	setting.Rule = richs
	fmt.Printf("%#v", setting)

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if err = dbusClient.AddZone(setting); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query.Setting)
}

// removeZone ...
// @Summary removeZone
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/port/remove [PUT]
func (this *SettingRouter) removeZone(c *gin.Context) {

	var query = &q.RemoveQuery{}
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

	if err = dbusClient.RemoveZone(query.Name); err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, query.Name)
}

// listZone ...
// @Summary listZone
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/setting/list [GET]
func (this *SettingRouter) listZone(c *gin.Context) {

	var query = &q.Query{}
	if err := c.BindQuery(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	zones, err := dbusClient.GetZones()
	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}
	q.SuccessResponse(c, code.OK, zones)
}

// defaultZone ...
// @Summary defaultZone
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/setting/default [GET]
func (this *SettingRouter) defaultZone(c *gin.Context) {

	var query = &q.Query{}
	if err := c.BindQuery(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	zone := dbusClient.GetDefaultZone()

	q.SuccessResponse(c, code.OK, zone)
}

// flushRuntime ...
// @Summary flushRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/setting/flushruntime [PUT]
func (this *SettingRouter) flush(c *gin.Context) {
	var query = &q.Query{}
	if err := c.BindQuery(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}
	if err := dbusClient.RuntimeFlush(query.Zone); err != nil {
		q.APIResponse(c, code.InternalServerError, err)
	}

	q.SuccessResponse(c, code.OK, query.Zone)
}

// setDefaultZone ...
// @Summary setDefaultZone
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v2/setting/default [PUT]
func (this *SettingRouter) setDefaultZone(c *gin.Context) {

	var query = &q.Query{}
	if err := c.BindQuery(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	if err := dbusClient.SetDefaultZone(query.Zone); err != nil {
		q.APIResponse(c, code.InternalServerError, err)
	}

	q.SuccessResponse(c, code.OK, query.Zone)
}
