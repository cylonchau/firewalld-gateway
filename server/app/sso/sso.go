package sso

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/praserx/ipconv"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	token2 "github.com/cylonchau/firewalld-gateway/utils/auther"
	userModel "github.com/cylonchau/firewalld-gateway/utils/model"
)

type SSO struct{}

// signinHandler godoc
// @Summary login.
// @Description login.
// @Tags SSO
// @Accept  json
// @Produce json
// @Param   query  body  query.UserQuery   false "signup body"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /sso/signin [post]
func (s *SSO) signinHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	userQuery := &query.UserQuery{}
	enconterError = c.ShouldBindJSON(&userQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	var (
		user         userModel.User
		token        string
		isPrivileged bool
		roles        []string
	)
	if user, enconterError = userModel.QueryUserWithUsername(userQuery.Username); enconterError == nil {
		if reflect.DeepEqual(user, userModel.User{}) {
			query.SuccessResponse(c, query.ErrUserNotExist, nil)
			return
		}

		if user.Username == userQuery.Username && user.Password == userModel.EncryptPassword(userQuery.Password) {
			var ip uint32
			switch user.ID {
			case 1:
				isPrivileged = true
				roles = []string{}
			default:
			}

			if ip, enconterError = userModel.GetRequestIP(c.Request); enconterError == nil {
				userModel.LastLogin(int64(user.ID), ip)
				if token, enconterError = token2.GenToken(int64(user.ID)); enconterError == nil {
					query.SuccessResponse(c, nil, query.UserResp{
						UserID:       uint64(user.ID),
						Token:        token,
						LoginIP:      ipconv.IntToIPv4(ip).String(),
						Roles:        roles,
						IsPrivileged: isPrivileged,
					})
					return
				}
			}

			if token, enconterError = token2.GenToken(int64(user.ID)); enconterError == nil {
				query.SuccessResponse(c, nil, query.UserResp{
					UserID:       uint64(user.ID),
					Token:        token,
					Roles:        roles,
					IsPrivileged: isPrivileged,
				})
				return
			}
		}
	}
	query.SuccessResponse(c, query.ErrPasswordIncorrect, nil)
}

// signupHandler godoc
// @Summary register.
// @Description register.
// @Tags SSO
// @Accept  json
// @Produce json
// @Param   query  body  query.UserQuery   false "user body"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /sso/signup [post]
func (s *SSO) signupHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	userQuery := &query.UserQuery{}
	var enconterError error
	enconterError = c.ShouldBindJSON(&userQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = userModel.CreateUser(userQuery); enconterError != nil {
		query.API409Response(c, enconterError)
		return
	}
	userQuery.Password = ""
	query.SuccessResponse(c, nil, userQuery)
}
