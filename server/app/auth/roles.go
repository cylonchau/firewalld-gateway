package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

// getRoles godoc
// @Summary Return all roles.
// @Description Return all roles.
// @Tags Auth
// @Accept  json
// @Produce json
// @Param   id  	query  int  false "token id"
// @Param   limit  	query  int  false "limit"
// @Param   offset  query  int  false "offset"
// @Param   sort  	query  string  false "sort"
// @Param   title  	query  string  false "sort"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles [get]
func (u *Auth) getRoles(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	rolesQuery := &query.ListQuery{}
	enconterError = c.Bind(&rolesQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if roles, enconterError := model.GetRoles(rolesQuery.Title, int(rolesQuery.Offset), int(rolesQuery.Limit), rolesQuery.Sort); enconterError == nil {
		if len(roles) > 0 {
			query.SuccessResponse(c, nil, roles)
			return
		}
		enconterError = errors.New(query.ErrRouterIsEmpty.Error())

	}
	query.SuccessResponse(c, enconterError, nil)
}

// getRoleByUserId godoc
// @Summary Return roles by user ID.
// @Description Return roles associated with a specific user ID.
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles/{id} [get]
func (u *Auth) getRoleByUserId(c *gin.Context) {
	var enconterError error
	rolesQuery := &query.QueryWithID{}
	if enconterError = c.ShouldBindUri(rolesQuery); enconterError != nil {
		query.API400Response(c, enconterError)
		return
	}
	// 3. 获取用户角色信息
	roles, err := model.GetRolesByUID(rolesQuery.ID)
	if err != nil {
		query.APIResponse(c, err, nil)
		return
	}

	// 4. 返回角色信息
	if len(roles) == 0 {
		query.SuccessResponse(c, nil, "No roles found for the user")
		return
	}

	query.SuccessResponse(c, nil, roles)
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
	roleQuery := &query.RoleEditQuery{}
	enconterError = c.ShouldBindJSON(&roleQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	routers := model.GenerateRouterWithID(roleQuery.RouterIDs)

	if enconterError = model.CreateRole(roleQuery, routers); enconterError != nil {
		query.SuccessResponse(c, enconterError, nil)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
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
	roleQuery := &query.RoleEditQuery{}
	enconterError = c.ShouldBindJSON(&roleQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if roleQuery.ID > 0 {
		if enconterError = model.UpdateRoleWithID(roleQuery); enconterError != nil {
			query.API409Response(c, enconterError)
			return
		}

		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}

// deleteRoleWithID godoc
// @Summary Delete a role with role_id.
// @Description Delete a role with role_id.
// @Tags Auth
// @Produce json
// @Accept x-www-form-urlencoded
// @Param id query query.IDQuery true "role id"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles [DELETE]
func (a *Auth) deleteRoleWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	roleQuery := &query.IDQuery{}
	enconterError = c.Bind(&roleQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteRoleWithID(roleQuery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
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
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/userRoles [get]
func (a *Auth) getUserRoles(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	roleQuery := &query.ListQuery{}
	enconterError = c.Bind(&roleQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if user, enconterError := model.GetRolesWithUID(roleQuery.ID); enconterError == nil {
		query.SuccessResponse(c, nil, user)
		return
	}
	query.SuccessResponse(c, enconterError, nil)
}

// getRoleRouters godoc
// @Summary Return all routers of role.
// @Description Return all routers of role。
// @Tags Auth
// @Accept  json
// @Produce json
// @Param id query int true "Input parameter"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roleRouters [get]
func (u *Auth) getRoleRouters(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	ids := strings.Split(c.Query("id"), ",")
	// 手动对请求参数进行详细的业务规则校验

	if routers, enconterError := model.GetRoutersWithRID(ids); enconterError == nil {
		query.SuccessResponse(c, nil, routers)
		return
	}
	query.SuccessResponse(c, enconterError, nil)
}
