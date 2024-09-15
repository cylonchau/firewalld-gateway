package v3

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

type NATRuleRouterV3 struct{}

func (this *NATRuleRouterV3) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/nat")
	portGroup.PUT("/", this.batchAddNATRuntime)
}

// batchAddNATRuntime godoc
// @Summary Add a NAT rule on firewalld runtime with delay timer.
// @Description Add a NAT rule on firewalld runtime with delay timer.
// @Tags firewalld NAT
// @Accept  json
// @Produce json
// @Param  query body  query.BatchForwardQuery false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/nat [put]
func (this *NATRuleRouterV3) batchAddNATRuntime(c *gin.Context) {
	var query = &api_query.BatchForwardQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	fmt.Println(query.Forwards[0].Forward)
	for _, item := range query.Forwards {
		go func(rule api_query.ForwardQuery) {
			contexts := context.WithValue(c, "action_obj", rule)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.CREATE_FORWARD)
			go batchFunction(contexts)
		}(item)
	}

	api_query.SuccessResponse(c, nil, api_query.BatchSuccessCreated)
}
