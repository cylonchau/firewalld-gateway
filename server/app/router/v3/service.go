package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/apis"
	code "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
)

type ServiceRouter struct{}

func (this *ServiceRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/service")
	portGroup.POST("/", this.batchAddService)
}

// reload ...
// @Summary reload
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v3/service [POST]
func (this *ServiceRouter) batchAddService(c *gin.Context) {
	var query = &apis.BatchServiceQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		apis.APIResponse(c, err, nil)
		return
	}
	for _, item := range query.Services {
		go func(p apis.ServiceQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "action_obj", p)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.CREATE_SERVICE)
			go batchFunction(contexts)
		}(item)
	}
	apis.SuccessResponse(c, code.OK, code.BatchSuccessCreated)
}
