package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

type SettingRouter struct{}

func (this *SettingRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/setting")
	portGroup.POST("/reload/runtime", this.reloadRuntime)
	portGroup.PUT("/sdzone", this.setDefautZone)
}

// reload godoc
// @Summary Reload firewalld runtime with delay timer.
// @Description Reload firewalld runtime with delay timer.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param  query body  query.BatchSettingQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/setting/reload/runtime [put]
func (this *SettingRouter) reloadRuntime(c *gin.Context) {
	var query = &api_query.BatchSettingQuery{}

	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}
	for _, item := range query.Hosts {
		go func(host string) {
			contexts := context.WithValue(c, "action_obj", item)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.RELOAD_FIREWALD)

			go batchFunction(contexts)
		}(item)
	}
	api_query.BacthMissionSuccessResponse(c, api_query.BatchSuccessCreated)
}

// batchAddPortRuntime godoc
// @Summary Set default zone on firewalld with delay timer, this is a runtime and permanent change.
// @Description Set default zone on firewalld with delay timer, this is a runtime and permanent change.
// @Tags firewalld setting
// @Accept  json
// @Produce json
// @Param  query body  query.BatchPortQuery  false "body"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/v3/setting/sdzone [POST]
func (this *SettingRouter) setDefautZone(c *gin.Context) {
	var query = &api_query.BatchZoneQuery{}

	if err := c.ShouldBindJSON(query); err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	for _, item := range query.ActionObject {
		go func(zoneAction api_query.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", zoneAction)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.SET_DEFAULT_ZONE)
			go batchFunction(contexts)
		}(item)
	}
	api_query.BacthMissionSuccessResponse(c, api_query.BatchSuccessCreated)
}
