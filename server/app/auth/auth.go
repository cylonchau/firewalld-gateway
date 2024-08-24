package auth

import (
	api_query "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	token2 "github.com/cylonchau/firewalld-gateway/utils/auther"
	userModel "github.com/cylonchau/firewalld-gateway/utils/model"

	"github.com/gin-gonic/gin"
)

type Auth struct{}

// allocateRolesToUser godoc
// @Summary allocate roles to user.
// @Description allocate roles to user.
// @Tags Auth
// @Produce json
// @Accept  json
// @Param userinput body query.RoleEditQuery true "Input parameter"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/roles/allocate [POST]
func (h *Auth) allocateRolesToUser(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &api_query.AllocateRouterQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		api_query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = userModel.AllocateRole2User(query.UserID, query.RoleIDs); enconterError != nil {
		api_query.API409Response(c, enconterError)
		return
	}

	api_query.SuccessResponse(c, api_query.OK, nil)
	return
}

// userInfoHandler godoc
// @Summary Return userinfo.
// @Description Return userinfo.
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string false "token"
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/auth/userinfo [get]
func (u *Auth) userInfoHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	userInfoQuery := &api_query.InfoQuery{}
	err := c.Bind(&userInfoQuery)

	// 手动对请求参数进行详细的业务规则校验
	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	uid, err := token2.GetInfo(userInfoQuery.Token)
	if err != nil {
		api_query.APIResponse(c, err, nil)
		return
	}

	var user userModel.User
	var enconterError error
	if user, enconterError = userModel.QueryUserWithUID(uid); enconterError == nil {
		api_query.SuccessResponse(c, nil, api_query.UserResp{
			Name: user.Username,
		})
		return
	}
	api_query.SuccessResponse(c, enconterError, nil)
}

// getClientIP godoc
// @Summary Get client login ip.
// @Description Get client login ip.
// @Tags Auth
// @Accept json
// @Produce json
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} interface{}
// @Router /security/auth/cip [get]
func (u *Auth) getClientIP(c *gin.Context) {
	api_query.SuccessResponse(c, nil, map[string]string{
		"ip": c.ClientIP(),
	})
}
