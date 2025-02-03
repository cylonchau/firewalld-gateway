package auth

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

// getRouters godoc
// @Summary Return all routers.
// @Description Return all routers.
// @Tags Auth
// @Accept  json
// @Produce json
// @Param   limit  	query  int   	false "limit"
// @Param   offset  query  int   	false "offset"
// @Param   sort  	query  string   false "sort"
// @Param   title  	query  string   false "sort"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/routers [get]
func (u *Auth) getRouters(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	routerQuery := &query.ListQuery{}
	enconterError = c.Bind(&routerQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if roles, enconterError := model.GetRouters(routerQuery.Title, int(routerQuery.Offset), int(routerQuery.Limit), routerQuery.Sort); enconterError == nil {
		if len(roles) > 0 {
			query.SuccessResponse(c, nil, roles)
			return
		}
		enconterError = errors.New(query.ErrRoleIsEmpty.Error())

	}
	query.SuccessResponse(c, enconterError, nil)
}

// getRoutersByRoleID godoc
// @Summary Return routers by role ID.
// @Description Return routers associated with a specific role ID.
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/routers/{id} [get]
func (u *Auth) getRoutersByRoleID(c *gin.Context) {
	var encounterError error
	roleQuery := &query.QueryWithID{}
	if encounterError = c.ShouldBindUri(roleQuery); encounterError != nil {
		query.API400Response(c, encounterError)
		return
	}

	routers, err := model.GetRoutersByRoleID(roleQuery.ID)
	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	// 4. 返回路由信息
	if len(routers) == 0 {
		query.SuccessResponse(c, nil, "No routers found for the role")
		return
	}

	query.SuccessResponse(c, nil, routers)
}
