package router

import (
	"github.com/gin-gonic/gin"

	code "github.com/cylonchau/firewalld-gateway/server/apis"
)

func ping(c *gin.Context) {
	code.SuccessResponse(c, code.OK, "pong")
}
