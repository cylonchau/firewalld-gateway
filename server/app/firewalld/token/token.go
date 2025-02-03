package Token

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

type Token struct{}

func (t *Token) RegisterTokenAPI(g *gin.RouterGroup) {
	g.GET("/", t.listToken)
	g.PUT("/", t.createToken)
	g.DELETE("/", t.deleteTokenWithID)
	g.POST("/", t.updateTokenWithID)
}

// createToken godoc
// @Summary Create a token.
// @Description Create a token.
// @Tags Token
// @Produce json
// @Accept json
// @Param userinput body query.TokenEditQuery true "token body"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/tokens [PUT]
func (t *Token) createToken(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tokenQuery := &query.TokenEditQuery{}
	enconterError = c.ShouldBindJSON(&tokenQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	if tokenQuery != nil {
		if enconterError = model.CreateToken(tokenQuery.SignedTo, tokenQuery.Description); enconterError != nil {
			query.APIResponse(c, enconterError, nil)
			return
		}
	} else {
		query.SuccessResponse(c, query.QUERY_NULL, nil)
	}
	query.SuccessResponse(c, query.OK, nil)
}

// listToken godoc
// @Summary Return all tokens.
// @Description Return all tokens.
// @Tags Token
// @Accept  json
// @Produce json
// @Param  id  query  int false "token id"
// @Param limit query  int false "limit"
// @Param offset query  int false "offset"
// @Param sort query  string false "sort"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/tokens [get]
func (t *Token) listToken(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tokenQuery := &query.ListQuery{}
	enconterError = c.Bind(&tokenQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	var list map[string]interface{}

	if list, enconterError = model.GetTokens(tokenQuery.Title, int(tokenQuery.Offset), int(tokenQuery.Limit), tokenQuery.Sort); enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, list)
}

// deleteTokenWithID godoc
// @Summary Delete a token with token_id.
// @Description Delete a token with token_id.
// @Tags Token
// @Produce json
// @Accept  x-www-form-urlencoded
// @Param id query query.IDQuery true "token id"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/tokens [DELETE]
func (t *Token) deleteTokenWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tokenQuery := &query.IDQuery{}
	enconterError = c.Bind(&tokenQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteTokenWithID(tokenQuery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// updateTokenWithID godoc
// @Summary Update token with token_id.
// @Description Update token with token_id.
// @Tags Token
// @Produce json
// @Accept  json
// @Param userinput body query.TokenEditQuery true "Input parameter"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /security/tokens [POST]
func (t *Token) updateTokenWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tokenQuery := &query.TokenEditQuery{}
	enconterError = c.ShouldBindJSON(&tokenQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if tokenQuery.ID > 0 {
		if enconterError = model.UpdateTokenWithID(tokenQuery.ID, tokenQuery.SignedTo, tokenQuery.Description, tokenQuery.IsUpdate); enconterError != nil {
			query.API409Response(c, enconterError)
			return
		}

		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}
