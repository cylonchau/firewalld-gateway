package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/apis"
	code "github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
)

type SettingRouter struct{}

func (this *SettingRouter) RegisterBatchAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/setting")
	portGroup.POST("/reload", this.reload)
	portGroup.PUT("/zone", this.setDefautZone)
}

// reload ...
// @Summary reload
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v3/setting/reload [POST]
func (this *SettingRouter) reload(c *gin.Context) {
	var query = &apis.BatchSettingQuery{}

	if err := c.ShouldBindJSON(query); err != nil {
		apis.APIResponse(c, err, nil)
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
	apis.BacthMissionSuccessResponse(c, code.BatchSuccessCreated)
}

// set default zone ...
// @Summary set default zone
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v3/setting/zone [POST]
func (this *SettingRouter) setDefautZone(c *gin.Context) {
	var query = &apis.BatchZoneQuery{}

	if err := c.ShouldBindJSON(query); err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	for _, item := range query.ActionObject {
		go func(zoneAction apis.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", zoneAction)
			contexts = context.WithValue(contexts, "delay_time", query.Delay)
			contexts = context.WithValue(contexts, "event_name", batch_processor.SET_DEFAULT_ZONE)
			go batchFunction(contexts)
		}(item)
	}
	apis.BacthMissionSuccessResponse(c, code.BatchSuccessCreated)
}
