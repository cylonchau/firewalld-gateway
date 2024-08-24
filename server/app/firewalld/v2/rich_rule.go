package v2

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type RichRuleV2Router struct{}

func (this *RichRuleV2Router) RegisterPortAPI(g *gin.RouterGroup) {
	richGroup := g.Group("/rich")
	richGroup.GET("/", this.getRichRulesOnPermanent)
	richGroup.PUT("/", this.addRichRuleOnPermanent)
	richGroup.DELETE("/", this.delRichRuleOnPermanent)
}

// getRichRulesOnPermanent godoc
// @Summary Get rich rule list on firewalld permanent.
// @Description Get rich rule list on firewalld permanent.
// @Tags firewalld rich
// @Accept json
// @Produce json
// @Param  ip  query  string true "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} []interface{}
// @Router /fw/v2/rich [get]
func (this *RichRuleV2Router) getRichRulesOnPermanent(c *gin.Context) {
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

	rules, err := dbusClient.GetPermanentRichRules(rich.Zone)

	if err != nil {
		api_query.APIResponse(c, err, rules)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, rules)
}

// addRichRuleOnPermanent godoc
// @Summary Add rich rule on firewalld permanent.
// @Description Add rich on firewalld permanent.
// @Tags firewalld rich
// @Accept json
// @Produce json
// @Param  query  body  query.RichQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/rich [put]
func (this *RichRuleV2Router) addRichRuleOnPermanent(c *gin.Context) {

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

	err = dbusClient.AddPermanentRichRule(query.Zone, query.Rich)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, query.Rich)
}

// delRichRuleOnPermanent godoc
// @Summary Remove rich rule on firewalld permanent.
// @Description Remove rich rule on firewalld permanent.
// @Tags firewalld rich
// @Accept json
// @Produce json
// @Param  query body query.RichQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v2/rich [delete]
func (this *RichRuleV2Router) delRichRuleOnPermanent(c *gin.Context) {

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

	err = dbusClient.RemovePermanentRichRule(query.Zone, query.Rich)

	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, query.Rich)
}
