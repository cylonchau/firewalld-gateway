package host

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	hostModel "github.com/cylonchau/firewalld-gateway/utils/model"
)

type Host struct{}

func (h *Host) RegisterHostAPI(g *gin.RouterGroup) {
	g.PUT("/", h.createHost)
	g.GET("/", h.listHost)
	g.POST("/", h.updateHostWithID)
	g.DELETE("/", h.deleteHostWithID)
}

// createHost godoc
// @Summary Add a new host into uranus.
// @Description Add a new host into uranus.
// @Tags Hosts
// @Accept  json
// @Produce json
// @Param query body query.HostQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/host [PUT]
func (h *Host) createHost(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	hostQuery := &query.HostQuery{}
	enconterError = c.ShouldBindJSON(&hostQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = hostModel.CreateHost(hostQuery.IP, hostQuery.Hostname, hostQuery.TagId); enconterError != nil {
		query.API409Response(c, enconterError)
		return
	}

	query.SuccessResponse(c, query.OK, nil)
}

// updateHostWithID godoc
// @Summary Update host information with host id.
// @Description Update host information with host id.
// @Tags Hosts
// @Accept  json
// @Produce json
// @Param query body query.HostQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/host [POST]
func (h *Host) updateHostWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	hostQuery := &query.HostQuery{}
	enconterError = c.ShouldBindJSON(&hostQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if hostQuery.ID > 0 {
		if enconterError = hostModel.UpdateHostWithID(hostQuery); enconterError != nil {
			query.API409Response(c, enconterError)
			return
		}

		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}

// listHost godoc
// @Summary Get host from uranus
// @Description Get host from uranus
// @Tags Hosts
// @Accept json
// @Produce json
// @Param query body query.ListHostQuery false "body"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /fw/host [GET]
func (h *Host) listHost(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	hostQuery := &query.ListHostQuery{}
	enconterError = c.Bind(&hostQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := hostModel.GetHosts(int(hostQuery.Offset), int(hostQuery.Limit), hostQuery.Sort)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}

	query.SuccessResponse(c, query.OK, list)
}

// deleteHostWithID godoc
// @Summary Delete host with host id.
// @Description Delete host with host id.
// @Tags Hosts
// @Accept  json
// @Produce json
// @Param query body query.IDQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/host [DELETE]
func (h *Host) deleteHostWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	hostQuery := &query.IDQuery{}
	enconterError = c.Bind(&hostQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = hostModel.DeleteHostWithID(hostQuery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}
