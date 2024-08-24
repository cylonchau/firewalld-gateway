package v1

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type RichRuleV1Router struct{}

func (this *RichRuleV1Router) RegisterPortAPI(g *gin.RouterGroup) {
	richGroup := g.Group("/rich")
	richGroup.PUT("/", this.addRichRuleAtRuntime)
	richGroup.GET("/", this.getRichRulesAtRuntime)
	richGroup.DELETE("/", this.delRichRuleAtRuntime)
}

// queryInRuntime godoc
// @Summary Get rich rule list at firewalld runtimes.
// @Description Get rich rule list at firewalld runtimes.
// @Tags firewalld rich
// @Accept json
// @Produce json
// @Param  ip  query  string true "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} []interface{}
// @Router /fw/v1/rich [get]
func (this *RichRuleV1Router) getRichRulesAtRuntime(c *gin.Context) {

	var rich = &api_query.Query{}

	if err := c.BindQuery(rich); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(rich.Ip)
	if err != nil {
		api_query.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	rules, err := dbusClient.GetRichRules(rich.Zone)

	if err != nil {
		api_query.APIResponse(c, err, rules)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, rules)
}

// addRichRuleAtRuntime godoc
// @Summary Add rich rule at firewalld runtimes.
// @Description Add rich rule at firewalld runtimes.
// @Tags firewalld rich
// @Accept json
// @Produce json
// @Param  query  body  query.RichQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/rich [put]
func (this *RichRuleV1Router) addRichRuleAtRuntime(c *gin.Context) {

	var query = &api_query.RichQuery{}

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

	err = dbusClient.AddRichRule(query.Zone, query.Rich, query.Timeout)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, query.Rich)
}

// delRichRuleAtRuntime godoc
// @Summary Remove rich rule at firewalld runtimes.
// @Description Remove rich rule at firewalld runtimes.
// @Tags firewalld rich
// @Accept json
// @Produce json
// @Param  query  body  query.RichQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v1/rich [delete]
func (this *RichRuleV1Router) delRichRuleAtRuntime(c *gin.Context) {

	var query = &api_query.RichQuery{}
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

	err = dbusClient.RemoveRichRule(query.Zone, query.Rich)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, query.Rich)
}
