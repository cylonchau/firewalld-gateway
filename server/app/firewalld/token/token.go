package Token

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

type Token struct{}

func (t *Token) RegisterTokenAPI(g *gin.RouterGroup) {
	g.GET("/", t.listToken)
	g.POST("/", t.createToken)
	g.DELETE("/", t.deleteTokenWithID)
	g.PUT("/", t.updateTokenWithID)
}

func (t *Token) createToken(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.TokenEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	if query != nil {
		if enconterError = model.CreateToken(query.SignedTo, query.Description); enconterError != nil {
			apis.APIResponse(c, enconterError, nil)
			return
		}
	} else {
		apis.SuccessResponse(c, apis.QUERY_NULL, nil)
	}
	apis.SuccessResponse(c, apis.OK, nil)
}

func (t *Token) listToken(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	var list map[string]interface{}

	if list, enconterError = model.GetTokens(int(query.Offset), int(query.Limit), query.Sort); enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, list)
}

func (t *Token) deleteTokenWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteTokenWithID(query.ID)
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, nil)
}

func (t *Token) updateTokenWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.TokenEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = model.UpdateTokenWithID(query.ID, query.SignedTo, query.Description, query.IsUpdate); enconterError != nil {
			apis.API409Response(c, enconterError)
			return
		}

		apis.SuccessResponse(c, apis.OK, nil)
		return
	}
	apis.APIResponse(c, errors.New("invaild id"), nil)
}
