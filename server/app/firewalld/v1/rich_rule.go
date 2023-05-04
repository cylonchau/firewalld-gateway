package v1

import (
	code "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type RichRuleRouter struct{}

func (this *RichRuleRouter) RegisterPortAPI(g *gin.RouterGroup) {
	richGroup := g.Group("/rich")
	richGroup.POST("/", this.addRichRuleAtRuntime)
	richGroup.GET("/", this.getRichRulesAtRuntime)
	richGroup.DELETE("/", this.delRichRuleAtRuntime)
}

// GetRichRules ...
// @Summary GetRichRules
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/rich/get [GET]
func (this *RichRuleRouter) getRichRulesAtRuntime(c *gin.Context) {

	var rich = &code.Query{}

	if err := c.BindQuery(rich); err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(rich.Ip)
	if err != nil {
		code.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	rules, err := dbusClient.GetRichRules(rich.Zone)

	if err != nil {
		code.APIResponse(c, err, rules)
		return
	}

	code.SuccessResponse(c, code.OK, rules)
}

// AddRichRule ...
// @Summary GetRichRules
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/rich/add [POST]
func (this *RichRuleRouter) addRichRuleAtRuntime(c *gin.Context) {

	var query = &code.RichQuery{}

	if err := c.BindJSON(query); err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		code.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	err = dbusClient.AddRichRule(query.Zone, query.Rich, query.Timeout)

	if err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	code.SuccessResponse(c, code.OK, query.Rich)
}

// DelRichRule ...
// @Summary DelRichRule
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/rich/delete [DELETE]
func (this *RichRuleRouter) delRichRuleAtRuntime(c *gin.Context) {

	var query = &code.RichQuery{}
	if err := c.BindJSON(query); err != nil {
		code.APIResponse(c, err, nil)
		return
	}
	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		code.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	err = dbusClient.RemoveRichRule(query.Zone, query.Rich)

	if err != nil {
		code.APIResponse(c, err, nil)
		return
	}

	code.SuccessResponse(c, code.OK, query.Rich)
}
