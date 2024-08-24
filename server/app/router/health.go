package router

import (
	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

// ping godoc
// @Summary Return process status.
// @Description Return process status.
// @Tags Health
// @Accept  json
// @Produce json
// @Success 200 {object} interface{}
// @Router /ping [get]
func ping(c *gin.Context) {
	query.SuccessResponse(c, query.OK, "pong")
}
