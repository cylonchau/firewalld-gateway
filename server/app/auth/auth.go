package auth

import (
	"reflect"

	"github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/server/app/auther"
	userModel "github.com/cylonchau/firewalld-gateway/utils/model"

	"github.com/gin-gonic/gin"
	"github.com/praserx/ipconv"
)

type Auth struct{}

func (a *Auth) RegisterUserAPI(g *gin.RouterGroup) {
	authGroup := g.Group("/")
	authGroup.POST("/signin", a.signinHandler)
	authGroup.POST("/signup", a.signupHandler)
	authGroup.GET("/info", a.userInfoHandler)
	authGroup.GET("/cip", a.getClientIP)

}

func (u *Auth) signinHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	userQuery := &apis.UserQuery{}
	enconterError = c.ShouldBindJSON(&userQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	var user userModel.User
	var token string
	if user, enconterError = userModel.QueryUserWithUsername(userQuery.Username); enconterError == nil {
		if reflect.DeepEqual(user, userModel.User{}) {
			apis.SuccessResponse(c, apis.ErrUserNotExist, nil)
			return
		}

		if user.Username == userQuery.Username && user.Password == userModel.EncryptPassword(userQuery.Password) {
			var ip uint32
			if ip, enconterError = userModel.GetRequestIP(c.Request); enconterError == nil {
				userModel.LastLogin(int64(user.ID), ip)
				if token, enconterError = auther.GenToken(int64(user.ID)); enconterError == nil {
					apis.SuccessResponse(c, nil, apis.UserResp{
						UserID:  uint64(user.ID),
						Token:   token,
						LoginIP: ipconv.IntToIPv4(ip).String(),
					})
					return
				}
			}

			if token, enconterError = auther.GenToken(int64(user.ID)); enconterError == nil {
				apis.SuccessResponse(c, nil, apis.UserResp{
					UserID: uint64(user.ID),
					Token:  token,
				})
				return
			}
		}
	}
	apis.SuccessResponse(c, apis.ErrPasswordIncorrect, nil)
}

func (u *Auth) signupHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	userQuery := &apis.UserQuery{}
	var enconterError error
	enconterError = c.ShouldBindJSON(&userQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = userModel.CreateUser(userQuery); enconterError != nil {
		apis.API409Response(c, enconterError)
		return
	}
	userQuery.Password = ""
	apis.SuccessResponse(c, nil, userQuery)
}

func (u *Auth) userInfoHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	userInfoQuery := &apis.InfoQuery{}
	err := c.Bind(&userInfoQuery)

	// 手动对请求参数进行详细的业务规则校验
	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	uid, err := auther.GetInfo(userInfoQuery.Token)
	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	var user userModel.User
	var enconterError error
	if user, enconterError = userModel.QueryUserWithUID(uid); enconterError == nil {
		apis.SuccessResponse(c, nil, apis.UserResp{
			Name: user.Username,
		})
		return
	}
	apis.SuccessResponse(c, enconterError, nil)
}

func (u *Auth) getClientIP(c *gin.Context) {
	apis.SuccessResponse(c, nil, map[string]string{
		"ip": c.ClientIP(),
	})
}
