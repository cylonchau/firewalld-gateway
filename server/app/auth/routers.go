package auth

import (
	"errors"

	"github.com/gin-gonic/gin"

	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

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
// @Router /security/auth/routers [get]
func (u *Auth) getRouters(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &api_query.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		api_query.APIResponse(c, enconterError, nil)
		return
	}

	if roles, enconterError := model.GetRouters(int(query.Offset), int(query.Limit), query.Sort); enconterError == nil {
		if len(roles) > 0 {
			api_query.SuccessResponse(c, nil, roles)
			return
		}
		enconterError = errors.New(api_query.ErrRoleIsEmpty.Error())

	}
	api_query.SuccessResponse(c, enconterError, nil)
}
