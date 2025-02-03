package audit

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

type Audit struct{}

func (a *Audit) RegisterAuditAPI(g *gin.RouterGroup) {
	// user
	userGroup := g.Group("/")
	userGroup.GET("/", a.getAuditLogs)
}

// getAuditLogs godoc
// @Summary Return all tokens.
// @Description Return all tokens.
// @Tags Audit
// @Accept  json
// @Produce json
// @Param   id  	query  int  	false "token id"
// @Param   limit  	query  int   	false "limit"
// @Param   offset  query  int   	false "offset"
// @Param   sort  	query  string   false "sort"
// @Param   title  	query  string   false "sort"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/audit [get]
func (a *Audit) getAuditLogs(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	auditQuery := &query.ListQuery{}
	enconterError = c.Bind(&auditQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if roles, enconterError := model.GetAuditLogs(auditQuery.Title, int(auditQuery.Offset), int(auditQuery.Limit), auditQuery.Sort); enconterError == nil {
		if len(roles) > 0 {
			query.SuccessResponse(c, nil, roles)
			return
		}
		enconterError = errors.New(query.ErrRoleIsEmpty.Error())

	}
	query.SuccessResponse(c, enconterError, nil)
}
