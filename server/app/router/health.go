package router

import (
	"github.com/gin-gonic/gin"

	q "github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
)

func ping(c *gin.Context) {
	q.SuccessResponse(c, code.OK, "pong")
}
