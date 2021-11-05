package v1

import (
	"firewall-api/code"
	"firewall-api/utils/dbus"
	q "firewall-api/utils/query"

	"github.com/gin-gonic/gin"
)

type RichRuleRouter struct{}

func (this *RichRuleRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/rich")
	portGroup.POST("/add", this.addRichRuleAtRuntime)
	portGroup.GET("/get", this.getRichRulesAtRuntime)
	portGroup.DELETE("/delete", this.delRichRuleAtRuntime)
}

// GetRichRules ...
// @Summary GetRichRules
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/rich/get [GET]
func (this *RichRuleRouter) getRichRulesAtRuntime(c *gin.Context) {

	var rich = &q.Query{}

	if err := c.BindQuery(rich); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := dbus.NewDbusClientService(rich.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	rules, err := dbusClient.GetRichRules(rich.Zone)

	if err != nil {
		q.APIResponse(c, err, rules)
		return
	}

	if len(rules) <= 0 {
		q.NotFount(c, code.ErrRichNotFount, rules)
		return
	}

	q.SuccessResponse(c, code.OK, rules)
}

// AddRichRule ...
// @Summary GetRichRules
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/rich/add [POST]
func (this *RichRuleRouter) addRichRuleAtRuntime(c *gin.Context) {

	var query = &q.RichQuery{}

	if err := c.BindJSON(query); err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	if query.Rich.Family == "" {
		query.Rich.Family = "ipv4"
	}

	dbusClient, err := dbus.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		q.ConnectDbusService(c, err)
		return
	}

	err = dbusClient.AddRichRule(query.Zone, query.Rich, query.Timeout)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	q.SuccessResponse(c, code.OK, query.Rich)
}

// DelRichRule ...
// @Summary DelRichRule
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/rich/delete [DELETE]
func (this *RichRuleRouter) delRichRuleAtRuntime(c *gin.Context) {

	var query = &q.RichQuery{}

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

	if query.Rich.Family == "" {
		query.Rich.Family = "ipv4"
	}

	err = dbusClient.RemoveRichRule(query.Zone, query.Rich)

	if err != nil {
		q.APIResponse(c, err, nil)
		return
	}

	q.SuccessResponse(c, code.OK, query.Rich)
}
