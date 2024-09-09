package v2

import (
	"strings"

	"github.com/cylonchau/firewalld-gateway/api"
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
	"github.com/ulule/deepcopier"
)

type SettingRouter struct{}

func (this *SettingRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/setting")
	portGroup.GET("/", this.listZone)
	portGroup.PUT("/", this.addZoneSetting)
	portGroup.DELETE("/", this.removeZone)
	portGroup.GET("/dz", this.defaultZone)
	portGroup.GET("/dp", this.defaultPolicy)
	portGroup.POST("/sdz", this.setDefaultZone)
	portGroup.POST("/reload", this.reload)
	portGroup.POST("/flush", this.flush)

}

// listZone godoc
// @Summary List zone.
// @Description List zone.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param  query  body  query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/setting [get]
func (this *SettingRouter) listZone(c *gin.Context) {

	var query = &api_query.Query{}
	if err := c.BindQuery(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	zones, err := dbusClient.GetZones()
	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, zones)
}

// defaultZone godoc
// @Summary Get default zone.
// @Description Get default zone.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param  query  body  query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/setting/dz [get]
func (this *SettingRouter) defaultZone(c *gin.Context) {

	var query = &api_query.Query{}
	if err := c.BindQuery(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	zone := dbusClient.GetDefaultZone()

	api_query.SuccessResponse(c, api_query.OK, zone)
}

// defaultPolicy godoc
// @Summary Get default policy.
// @Description Get default policy.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param  query  body  query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/setting/dp [get]
func (this *SettingRouter) defaultPolicy(c *gin.Context) {
	var query = &api_query.Query{}
	if err := c.BindQuery(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	zone := dbusClient.GetDefaultPolicy()

	api_query.SuccessResponse(c, api_query.OK, zone)
}

// reload godoc
// @Summary Reload firewalld.
// @Description Reload firewalld.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param ip query string true "ip"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/setting/reload [post]
func (this *SettingRouter) reload(c *gin.Context) {

	var query = &api_query.Query{}
	if err := c.Bind(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.Reload(); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, nil)
}

// setDefaultZone godoc
// @Summary Set default zone in firewalld.
// @Description Set default zone in firewalld.
// @Tags firewalld setting
// @Accept json
// @Produce json
// @Param ip query string true "ip"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/setting/sdz [post]
func (this *SettingRouter) setDefaultZone(c *gin.Context) {

	var query = &api_query.Query{}
	if err := c.BindQuery(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err := dbusClient.SetDefaultZone(query.Zone); err != nil {
		if strings.Contains(err.Error(), "INVALID_ZONE") {
			api_query.NotFount(c, api_query.ErrZoneNotFount, err)
			return
		}
		api_query.APIResponse(c, api_query.InternalServerError, err)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, query.Zone)
}

// flush godoc
// @Summary Flush all firewalld rules to default.
// @Description Flush all firewalld rules to default.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param query body query.Query  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/setting/flush [post]
func (this *SettingRouter) flush(c *gin.Context) {
	var query = &api_query.Query{}
	if err := c.BindQuery(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	if query.Zone == "" {
		query.Zone = "public"
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()
	if err := dbusClient.RuntimeFlush(query.Zone); err != nil {
		api_query.APIResponse(c, api_query.InternalServerError, err)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, query.Zone)
}

// addZoneSetting godoc
// @Summary Add setting rule in firewalld.
// @Description Add setting rule in firewalld.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param query  body  query.ZoneSettingQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/setting [put]
func (this *SettingRouter) addZoneSetting(c *gin.Context) {
	var query = &api_query.ZoneSettingQuery{}
	if err := c.BindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	var setting = &api.Settings{}
	var richs []string
	for _, n := range query.Setting.Rule {
		richs = append(richs, n.ToString())
	}

	deepcopier.Copy(query.Setting).To(setting)
	setting.Rule = richs

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.AddZone(setting); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query.Setting)
}

// removeZone godoc
// @Summary Remove a setting rules in firewalld.
// @Description Remove a setting rules in firewalld.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param  query  body  query.RemoveQuery   false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /fw/v2/setting [delete]
func (this *SettingRouter) removeZone(c *gin.Context) {

	var query = &api_query.RemoveQuery{}
	if err := c.BindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err = dbusClient.RemoveZone(query.Name); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, query.Name)
}
