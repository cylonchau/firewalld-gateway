package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

type ServiceRouter struct{}

func (this *ServiceRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/service")
	portGroup.PUT("/", this.batchAddServiceRuntime)
}

// batchAddServiceRuntime godoc
// @Summary Add a service rule on firewalld runtime with delay timer.
// @Description Add a service rule on firewalld runtime with delay timer.
// @Tags firewalld service
// @Accept json
// @Produce json
// @Param query body query.BatchServiceQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/service [put]
func (this *ServiceRouter) batchAddServiceRuntime(c *gin.Context) {
	var query = &api_query.BatchServiceQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	for _, item := range query.Services {
		go func(p api_query.ServiceQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "action_obj", p)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.CREATE_SERVICE)
			go batchFunction(contexts)
		}(item)
	}
	api_query.SuccessResponse(c, api_query.OK, api_query.BatchSuccessCreated)
}
