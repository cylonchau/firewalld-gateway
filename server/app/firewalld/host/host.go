package host

import (
	"github.com/gin-gonic/gin"

	hostModel "github.com/cylonchau/firewalld-gateway/model"
	"github.com/cylonchau/firewalld-gateway/server/apis"
)

type Host struct{}

func (h *Host) RegisterHostAPI(g *gin.RouterGroup) {
	g.POST("/", h.createHost)
	g.GET("/", h.listHost)
	g.DELETE("/", h.deleteHostWithID)
}

func (h *Host) createHost(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.HostQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = hostModel.CreateHost(query); enconterError != nil {
		apis.API409Response(c, enconterError)
		return
	}

	apis.SuccessResponse(c, apis.OK, nil)
}

func (h *Host) listHost(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.ListHostQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := hostModel.GetHosts(int(query.Offset), int(query.Limit))
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, list)
}

func (h *Host) deleteHostWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = hostModel.DeleteHostWithID(query.ID)
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, nil)
}
