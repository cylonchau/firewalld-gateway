package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

type MasqueradeRouter struct{}

func (this *MasqueradeRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/nat")
	portGroup.PUT("/", this.batchEnableMasquerade)
	portGroup.DELETE("/", this.batchDisableMasquerade)
}

// batchEnableMasquerade godoc
// @Summary Enable masqerade on firewalld runtime with delay timer.
// @Description Enable masqerade on firewalld runtime with delay timer.
// @Tags firewalld masquerade
// @Accept  json
// @Produce json
// @Param  query body  query.BatchZoneQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/masquerade [put]
func (this *MasqueradeRouter) batchEnableMasquerade(c *gin.Context) {
	var query = &api_query.BatchZoneQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	for _, item := range query.ActionObject {
		go func(host api_query.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", host)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.ENABLE_MASQUERADE)
			go batchFunction(contexts)
		}(item)
	}

	api_query.SuccessResponse(c, nil, api_query.BatchSuccessCreated)
}

// batchDisableMasquerade godoc
// @Summary Disable masqerade on firewalld runtime with delay timer.
// @Description Disable masqerade on firewalld runtime with delay timer.
// @Tags firewalld masquerade
// @Accept  json
// @Produce json
// @Param  query body  query.BatchZoneQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/masquerade [delete]
func (this *MasqueradeRouter) batchDisableMasquerade(c *gin.Context) {
	var query = &api_query.BatchZoneQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	for _, item := range query.ActionObject {
		go func(host api_query.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", host)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.DISABLE_MASQUERADE)
			go batchFunction(contexts)
		}(item)
	}

	api_query.SuccessResponse(c, nil, api_query.BatchSuccessCreated)
}
