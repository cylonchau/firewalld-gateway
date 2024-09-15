package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

type PortRouter struct{}

func (this *PortRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/ports")
	portGroup.PUT("/", this.batchAddPortRuntime)
	portGroup.PUT("/permanent", this.batchAddPortPerment)
}

// batchAddPortRuntime godoc
// @Summary Add a port rule on firewalld runtime with delay timer.
// @Description Add a port rule on firewalld runtime with delay timer.
// @Tags firewalld port
// @Accept  json
// @Produce json
// @Param  query body  query.BatchPortQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/ports [put]
func (this *PortRouter) batchAddPortRuntime(c *gin.Context) {
	var query = &api_query.BatchPortQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	for _, item := range query.Ports {
		go func(p api_query.PortQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "action_obj", p)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.CREATE_PORT)
			go batchFunction(contexts)
		}(item)
	}
	api_query.SuccessResponse(c, api_query.OK, api_query.BatchSuccessCreated)
}

// batchAddPortPerment godoc
// @Summary Add a port rule on firewalld permanent with delay timer.
// @Description Add a port rule on firewalld permanent with delay timer.
// @Tags firewalld port
// @Accept  json
// @Produce json
// @Param  query body  query.BatchPortQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/ports/permanent [put]
func (this *PortRouter) batchAddPortPerment(c *gin.Context) {
	var query = &api_query.BatchPortQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	for _, item := range query.Ports {
		go func(p api_query.PortQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "action_obj", p)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.CREATE_PORT_PERMANENT)
			go batchFunction(contexts)
		}(item)
	}
	api_query.SuccessResponse(c, api_query.OK, api_query.BatchSuccessCreated)
}
