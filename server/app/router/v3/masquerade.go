package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/apis"
	code_api "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
)

type MasqueradeRouter struct{}

func (this *MasqueradeRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/nat")
	portGroup.PUT("/", this.batchEnableMasquerade)
	portGroup.DELETE("/", this.batchDisableMasquerade)
}

// enable masquerade ...
// @Summary enable masquerade
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v3/nat [POST]
func (this *MasqueradeRouter) batchEnableMasquerade(c *gin.Context) {
	var query = &apis.BatchZoneQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	for _, item := range query.ActionObject {
		go func(host apis.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", host)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.ENABLE_MASQUERADE)
			go batchFunction(contexts)
		}(item)
	}

	apis.SuccessResponse(c, code_api.OK, code_api.BatchSuccessCreated)
}

// disable masquerade ...
// @Summary disable masquerade
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v3/nat [POST]
func (this *MasqueradeRouter) batchDisableMasquerade(c *gin.Context) {
	var query = &apis.BatchZoneQuery{}
	if err := c.ShouldBindJSON(query); err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	for _, item := range query.ActionObject {
		go func(host apis.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", host)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.DISABLE_MASQUERADE)
			go batchFunction(contexts)
		}(item)
	}

	apis.SuccessResponse(c, code_api.OK, code_api.BatchSuccessCreated)
}
