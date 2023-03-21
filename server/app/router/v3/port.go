package v3

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/server/batch_processor"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"
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
		function := func(c context.Context) {
			b := c.Value("portRule")
			port := b.(apis.PortQuery)
			dbusClient, err := firewalld.NewDbusClientService(port.Ip)
			if err != nil {
				return
			}
			defer func() {
				c.Done()
				dbusClient.Destroy()
			}()
			tName := batch_processor.RandName()
			event := batch_processor.Event{
				EventName: batch_processor.CREATE_PORT,
				Host:      port.Ip,
				TaskName:  tName,
				Task:      port,
			}
			batch_processor.P.Add(tName, event)

		}
		go func(p apis.PortQuery) {
			contexts := context.TODO()
			contexts = context.WithValue(contexts, "portRule", p)
			go function(contexts)
		}(item)
	}
	apis.SuccessResponse(c, code.OK, code.BatchSuccessCreated)
}
