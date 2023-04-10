package v1

import (
	code "github.com/cylonchau/firewalld-gateway/server/apis"

	"github.com/gin-gonic/gin"
)

type StoageRouter struct{}

func (this *StoageRouter) RegisterStoageAPI(g *gin.RouterGroup) {
	storeGroup := g.Group("/stroage")
	storeGroup.GET("/get", this.get)
}

func (this *StoageRouter) get(c *gin.Context) {
	code.SuccessResponse(c, code.OK, "")
}