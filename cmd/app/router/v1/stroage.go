package v1

import (
	"firewall-api/code"
	q "firewall-api/utils/query"

	"github.com/gin-gonic/gin"
)

type StoageRouter struct{}

func (this *StoageRouter) RegisterStoageAPI(g *gin.RouterGroup) {
	storeGroup := g.Group("/stoage")
	storeGroup.GET("/get", this.get)
}

func (this *StoageRouter) get(c *gin.Context) {
	q.SuccessResponse(c, code.OK, "")
}
