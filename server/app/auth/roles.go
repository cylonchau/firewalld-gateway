package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

// getRoles godoc
// @Summary Return all roles.
// @Description Return all roles.
// @Tags Auth
// @Accept  json
// @Produce json
// @Param   id  query  int   false "token id"
// @Param   limit  query  int   false "limit"
// @Param   offset  query  int   false "offset"
// @Param   sort  query  string   false "sort"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles [get]
func (u *Auth) getRoles(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &api_query.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		api_query.APIResponse(c, enconterError, nil)
		return
	}

	if roles, enconterError := model.GetRoles(int(query.Offset), int(query.Limit), query.Sort); enconterError == nil {
		if len(roles) > 0 {
			api_query.SuccessResponse(c, nil, roles)
			return
		}
		enconterError = errors.New(api_query.ErrRouterIsEmpty.Error())

	}
	api_query.SuccessResponse(c, enconterError, nil)
}

// createRole godoc
// @Summary Create a role.
// @Description Create a role.
// @Tags Auth
// @Produce json
// @Accept  json
// @Param userinput body query.RoleEditQuery true "role body"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles [PUT]
func (a *Auth) createRole(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	roleQuery := &api_query.RoleEditQuery{}
	enconterError = c.ShouldBindJSON(&roleQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		api_query.APIResponse(c, enconterError, nil)
		return
	}

	routers := model.GenerateRouterWithID(roleQuery.RouterIDs)

	if enconterError = model.CreateRole(roleQuery, routers); enconterError != nil {
		api_query.SuccessResponse(c, enconterError, nil)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, nil)
}

// updateRole godoc
// @Summary Update token with role_id.
// @Description Update token with role_id.
// @Tags Auth
// @Produce json
// @Accept  json
// @Param userinput body query.RoleEditQuery true "Input parameter"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles [POST]
func (a *Auth) updateRole(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &api_query.RoleEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		api_query.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = model.UpdateRoleWithID(query); enconterError != nil {
			api_query.API409Response(c, enconterError)
			return
		}

		api_query.SuccessResponse(c, api_query.OK, nil)
		return
	}
	api_query.APIResponse(c, errors.New("invaild id"), nil)
}

// deleteRoleWithID godoc
// @Summary Delete a role with role_id.
// @Description Delete a role with role_id.
// @Tags Auth
// @Produce json
// @Accept  x-www-form-urlencoded
// @Param id query query.IDQuery true "role id"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles [DELETE]
func (a *Auth) deleteRoleWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &api_query.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		api_query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteRoleWithID(query.ID)
	if enconterError != nil {
		api_query.API500Response(c, enconterError)
		return
	}
	api_query.SuccessResponse(c, api_query.OK, nil)
}

// getRouters godoc
// @Summary Return all routers.
// @Description Return all routers.
// @Tags Auth
// @Accept  json
// @Produce json
// @Param   id  query  int   false "token id"
// @Param   limit  query  int   false "limit"
// @Param   offset  query  int   false "offset"
// @Param   sort  query  string   false "sort"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/userRoles [get]
func (a *Auth) getUserRoles(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &api_query.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		api_query.APIResponse(c, enconterError, nil)
		return
	}

	if user, enconterError := model.GetRolesWithUID(query.ID); enconterError == nil {
		api_query.SuccessResponse(c, nil, user)
		return
	}
	api_query.SuccessResponse(c, enconterError, nil)
}

// getRoleRouters godoc
// @Summary Return all routers of role.
// @Description Return all routers of role。
// @Tags Auth
// @Accept  json
// @Produce json
// @Param id query int true "Input parameter"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roleRouters [get]
func (u *Auth) getRoleRouters(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	ids := strings.Split(c.Query("id"), ",")
	// 手动对请求参数进行详细的业务规则校验

	if routers, enconterError := model.GetRoutersWithRID(ids); enconterError == nil {
		api_query.SuccessResponse(c, nil, routers)
		return
	}
	api_query.SuccessResponse(c, enconterError, nil)
}
