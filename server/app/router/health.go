package router

import (
	"github.com/gin-gonic/gin"

	q "github.com/cylonchau/firewalld-gateway/apis"
	code "github.com/cylonchau/firewalld-gateway/server/apis"
)

func ping(c *gin.Context) {
	q.SuccessResponse(c, code.OK, "pong")
}
