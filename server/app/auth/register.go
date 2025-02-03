package auth

import "github.com/gin-gonic/gin"

func (a *Auth) RegisterAuthAPI(g *gin.RouterGroup) {
	// user
	authGroup := g.Group("/")
	authGroup.GET("/userinfo", a.userInfoHandler)
	authGroup.GET("/cip", a.getClientIP)

	// router handler
	authGroup.GET("/routers", a.getRouters)
	authGroup.GET("/routers/:id", a.getRoutersByRoleID)
	// role
	authGroup.GET("/roles", a.getRoles)
	authGroup.GET("/roles/:id", a.getRoleByUserId)
	authGroup.PUT("/roles", a.createRole)
	authGroup.POST("/roles", a.updateRole)
	authGroup.DELETE("/roles", a.deleteRoleWithID)
	authGroup.POST("/roles/allocate", a.allocateRolesToUser)
	// user
	authGroup.GET("/userRoles", a.getUserRoles)
	authGroup.GET("/roleRouters", a.getRoleRouters)

}
