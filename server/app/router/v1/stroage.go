package v1

import (
	q "github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"

	"github.com/gin-gonic/gin"
)

type StoageRouter struct{}

func (this *StoageRouter) RegisterStoageAPI(g *gin.RouterGroup) {
	storeGroup := g.Group("/stroage")
	storeGroup.GET("/get", this.get)
}

func (this *StoageRouter) get(c *gin.Context) {
	q.SuccessResponse(c, code.OK, "")
}
