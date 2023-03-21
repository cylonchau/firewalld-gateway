package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/server/batch_processor"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"
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
		function := func(c context.Context) {
			b := c.Value("service")
			service := b.(apis.ServiceQuery)
			dbusClient, err := firewalld.NewDbusClientService(service.Ip)
			if err != nil {
				return
			}
			defer func() {
				c.Done()
				dbusClient.Destroy()
			}()
			tName := batch_processor.RandName()
			event := batch_processor.Event{
				EventName: batch_processor.CREATE_SERVICE,
				Host:      service.Ip,
				TaskName:  tName,
				Task:      service,
			}
			batch_processor.P.Add(tName, event)

		}
		go func(p apis.ServiceQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "service", p)
			go function(contexts)
		}(item)
	}
	apis.SuccessResponse(c, code.OK, code.BatchSuccessCreated)
}
