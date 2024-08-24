package sso

import "github.com/gin-gonic/gin"

func (s *SSO) RegisterUserAPI(g *gin.RouterGroup) {
	// user
	authGroup := g.Group("/")
	authGroup.POST("/signin", s.signinHandler)
	authGroup.POST("/signup", s.signupHandler)
}
