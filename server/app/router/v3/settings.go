package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/server/batch_processor"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"
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
		function := func(c context.Context) {
			b := c.Value("settingQuery")
			host := b.(string)
			dbusClient, err := firewalld.NewDbusClientService(host)
			if err != nil {
				return
			}
			defer func() {
				c.Done()
				dbusClient.Destroy()
			}()
			tName := batch_processor.RandName()
			event := batch_processor.Event{
				EventName: batch_processor.RELOAD_FIREWALD,
				Host:      host,
				TaskName:  tName,
			}
			batch_processor.P.Add(tName, event)

		}
		go func(host string) {
			contexts := context.WithValue(c, "settingQuery", host)
			go function(contexts)
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
		function := func(c context.Context) {
			b := c.Value("setDefautZone")
			action := b.(apis.ZoneDst)
			dbusClient, err := firewalld.NewDbusClientService(action.Host)
			if err != nil {
				return
			}
			defer func() {
				c.Done()
				dbusClient.Destroy()
			}()
			tName := batch_processor.RandName()
			event := batch_processor.Event{
				EventName: batch_processor.SET_DEFAULT_ZONE,
				Host:      action.Host,
				TaskName:  tName,
				Task:      action.Zone,
			}
			batch_processor.P.Add(tName, event)

		}
		go func(zoneAction apis.ZoneDst) {
			contexts := context.WithValue(c, "setDefautZone", zoneAction)
			go function(contexts)
		}(item)
	}
	apis.BacthMissionSuccessResponse(c, code.BatchSuccessCreated)
}
