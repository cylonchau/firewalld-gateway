package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalldGateway/apis"
	code_api "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/server/batch_processor"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"
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
		function := func(c context.Context) {
			b := c.Value("action_obj")
			obj := b.(apis.ZoneDst)
			dbusClient, err := firewalld.NewDbusClientService(obj.Host)
			if err != nil {
				return
			}
			defer func() {
				c.Done()
				dbusClient.Destroy()
			}()
			tName := batch_processor.RandName()
			event := batch_processor.Event{
				EventName: batch_processor.ENABLE_MASQUERADE,
				Host:      obj.Host,
				TaskName:  tName,
				Task:      obj.Zone,
			}
			batch_processor.P.Add(tName, event)

		}
		go func(host apis.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", host)
			go function(contexts)
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
		function := func(c context.Context) {
			b := c.Value("action_obj")
			obj := b.(apis.ZoneDst)
			dbusClient, err := firewalld.NewDbusClientService(obj.Host)
			if err != nil {
				return
			}
			defer func() {
				c.Done()
				dbusClient.Destroy()
			}()
			tName := batch_processor.RandName()
			event := batch_processor.Event{
				EventName: batch_processor.DISABLE_MASQUERADE,
				Host:      obj.Host,
				TaskName:  tName,
				Task:      obj.Zone,
			}
			batch_processor.P.Add(tName, event)

		}
		go func(host apis.ZoneDst) {
			contexts := context.WithValue(c, "action_obj", host)
			go function(contexts)
		}(item)
	}

	apis.SuccessResponse(c, code_api.OK, code_api.BatchSuccessCreated)
}
