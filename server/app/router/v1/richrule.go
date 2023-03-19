package v1

import (
	q "github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"

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

	dbusClient, err := firewalld.NewDbusClientService(rich.Ip)
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

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
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

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
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
