package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/apis"
	code "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
)

type PortRouter struct{}

func (this *PortRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/port")
	portGroup.POST("/", this.batchAddPort)
}

// reload ...
// @Summary reload
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v3/port/add [POST]
func (this *PortRouter) batchAddPort(c *gin.Context) {
	var query = &apis.BatchPortQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		apis.APIResponse(c, err, nil)
		return
	}
	for _, item := range query.Ports {
		go func(p apis.PortQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "action_obj", p)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.CREATE_PORT)
			go batchFunction(contexts)
		}(item)
	}
	apis.SuccessResponse(c, code.OK, code.BatchSuccessCreated)
}
