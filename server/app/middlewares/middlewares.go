package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/utils/auther"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization： Bearer xxxxxx.xxxx.xxx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			apis.AuthFailed(c, apis.ErrNeedAuth, nil)
			c.Abort() // 中止
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			apis.Auth403Failed(c, apis.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		tokenStr := parts[1]
		if model.TokenIsDestoryed(tokenStr) {
			apis.Auth403Failed(c, apis.ErrTokenDestoryed, nil)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := auther.ParseToken(tokenStr)
		if err != nil {
			apis.Auth403Failed(c, apis.ErrTokenInvalid, err.Error())
			c.Abort()
			return
		}
		// 将当前请求的userid信息保存到请求的上下文c上
		c.Set(auther.UserIDKey, mc.UserID)

		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
