package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

type RichRuleRouterV3 struct{}

func (this *RichRuleRouterV3) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/rich")
	portGroup.PUT("/", this.batchAddRichRuntime)
}

// batchAddRichRuntime godoc
// @Summary Add a rich rules on firewalld runtime with delay timer.
// @Description Add a rich rules on firewalld runtime with delay timer.
// @Tags firewalld rich
// @Accept json
// @Produce json
// @Param query body query.BatchRichQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/rich [put]
func (this *RichRuleRouterV3) batchAddRichRuntime(c *gin.Context) {
	var batchRichQuery = &query.BatchRichQuery{}
	if err := c.ShouldBindJSON(batchRichQuery); err != nil {
		query.APIResponse(c, err, nil)
		return
	}
	for _, item := range batchRichQuery.Richs {
		go func(p query.RichQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "action_obj", p)
			contexts = context.WithValue(contexts, "delay_time", batchRichQuery.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.CREATE_RICH)
			go batchFunction(contexts)
		}(item)
	}
	query.SuccessResponse(c, query.OK, query.BatchSuccessCreated)
}
