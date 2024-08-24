package auth

import "github.com/gin-gonic/gin"

func (a *Auth) RegisterUserAPI(g *gin.RouterGroup) {
	// user
	authGroup := g.Group("/")
	authGroup.GET("/userinfo", a.userInfoHandler)
	authGroup.GET("/cip", a.getClientIP)

	// router handler
	authGroup.GET("/routers", a.getRouters)
	// role
	authGroup.GET("/roles", a.getRoles)
	authGroup.PUT("/roles", a.createRole)
	authGroup.POST("/roles", a.updateRole)
	authGroup.DELETE("/roles", a.deleteRoleWithID)
	authGroup.POST("/roles/allocate", a.allocateRolesToUser)
	// user
	authGroup.GET("/userRoles", a.getUserRoles)
	authGroup.GET("/roleRouters", a.getRoleRouters)
}
