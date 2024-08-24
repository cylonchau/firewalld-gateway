package user

import (
	"errors"
	"strings"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"

	"github.com/gin-gonic/gin"
)

type User struct{}

func (u *User) RegisterUserAPI(g *gin.RouterGroup) {
	// user
	userGroup := g.Group("/")
	userGroup.POST("/", u.updateUser)
	userGroup.DELETE("/", u.deleteUserWithID)
	userGroup.GET("/", u.getUsers)
	userGroup.PUT("/", u.createUser)
	userGroup.POST("/allocate", u.allocateRolesToUser)

}

// getUsers godoc
// @Summary Return all users.
// @Description Return all users。
// @Tags Users
// @Accept json
// @Produce json
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/users [get]
func (u *User) getUsers(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	q := &query.ListQuery{}
	enconterError = c.Bind(&q)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if roles, enconterError := model.GetUsers(int(q.Offset), int(q.Limit), q.Sort); enconterError == nil {
		if len(roles) > 0 {
			query.SuccessResponse(c, nil, roles)
			return
		}
		enconterError = errors.New(query.ErrRouterIsEmpty.Error())

	}
	query.SuccessResponse(c, enconterError, nil)
}

// createUser godoc
// @Summary Create a user.
// @Description Create a user。
// @Tags Users
// @Accept  json
// @Produce json
// @Param userinput body query.UserQuery true "Input parameter"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/users [put]
func (u *User) createUser(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	userQuery := &query.UserQuery{}
	enconterError = c.ShouldBindJSON(&userQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = model.CreateUser(userQuery); enconterError != nil {
		query.SuccessResponse(c, enconterError, nil)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// updateUser godoc
// @Summary update user information.
// @Description Create a user。
// @Tags Users
// @Produce json
// @Accept  json
// @Param userinput body query.UserEditQuery true "Input parameter"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/users [post]
func (u *User) updateUser(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	q := &query.UserEditQuery{}
	enconterError = c.ShouldBindJSON(&q)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if q.ID > 0 {
		if enconterError = model.UpdateUserWithID(q); enconterError != nil {
			query.API409Response(c, enconterError)
			return
		}

		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}

// deleteRoleWithID godoc
// @Summary Delete a user.
// @Description Delete a user.
// @Tags Users
// @Produce json
// @Accept  x-www-form-urlencoded
// @Param id query query.IDQuery true "Input parameter"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/users [delete]
func (u *User) deleteRoleWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	q := &query.IDQuery{}
	enconterError = c.Bind(&q)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteRoleWithID(q.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

func (u *User) getUserRoles(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	q := &query.ListQuery{}
	enconterError = c.Bind(&q)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if user, enconterError := model.GetRolesWithUID(q.ID); enconterError == nil {
		query.SuccessResponse(c, nil, user)
		return
	}
	query.SuccessResponse(c, enconterError, nil)
}

func (u *User) getRoleRouters(c *gin.Context) {
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

func (u *User) getUserWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	q := &query.IDQuery{}
	enconterError = c.Bind(&q)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteUserWithID(q.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

func (u *User) deleteUserWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	q := &query.IDQuery{}
	enconterError = c.Bind(&q)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteUserWithID(q.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// updateUser allocateRolesToUser
// @Summary Assign roles to users.
// @Description Assign roles to users。
// @Tags Users
// @Produce json
// @Accept  json
// @Param id body query.AllocateRoleQuery true "Input parameter"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/users/allocate [post]
func (h *User) allocateRolesToUser(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	q := &query.AllocateRoleQuery{}
	enconterError = c.ShouldBindJSON(&q)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = model.AllocateRole2User(q.UserID, q.RoleIDs); enconterError != nil {
		query.API409Response(c, enconterError)
		return
	}

	query.SuccessResponse(c, query.OK, nil)
	return
}
