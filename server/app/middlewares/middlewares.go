package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/auther"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

func writeLog(r *http.Request, id int64) {
	ip, _ := model.GetRequestIP(r)
	ua := user_agent.New(r.UserAgent())
	browserName, _ := ua.Browser()
	osName := ua.OS()

	auditLogData := map[string]interface{}{
		"user_id": id,
		"ip":      ip,
		"method":  r.Method,
		"path":    r.URL.Path,
		"browser": browserName,
		"system":  osName,
	}
	model.AppendAuditLog(auditLogData)
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization： Bearer xxxxxx.xxxx.xxx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			query.AuthFailed(c, query.ErrNeedAuth, nil)
			c.Abort() // 中止
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			query.Auth403Failed(c, query.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		tokenStr := parts[1]
		if model.TokenIsDestoryed(tokenStr) {
			query.Auth403Failed(c, query.ErrTokenDestoryed, nil)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := auther.ParseToken(tokenStr)
		if err != nil {
			query.Auth403Failed(c, query.ErrTokenInvalid, err.Error())
			c.Abort()
			return
		}

		if mc.UserID == 1 {
			c.Set(auther.UserIDKey, mc.UserID)
			writeLog(c.Request, mc.UserID)
			c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
			return
		}

		roles, err := model.GetRolesWithUID(int(mc.UserID))
		if err != nil {
			query.Auth403Failed(c, query.ErrTokenInvalid, err.Error())
			c.Abort()
			return
		}

		var ids []string
		for _, v := range roles.Roles {
			ids = append(ids, strconv.Itoa(v.ID))
		}
		routers, err := model.GetRoutersWithRID(ids)

		for _, v := range routers {
			if v.Path == c.Request.URL.Path && v.Method == c.Request.Method {
				// 将当前请求的userid信息保存到请求的上下文c上
				c.Set(auther.UserIDKey, mc.UserID)
				writeLog(c.Request, mc.UserID)
				c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
				return
			}
		}
		query.AuthNoPermission(c, query.ErrNoPermission)
		c.Abort()
		return
	}
}
