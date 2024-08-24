package Token

import (
	"errors"

	"github.com/gin-gonic/gin"

	query2 "github.com/cylonchau/firewalld-gateway/utils/apis/query"
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
	query := &query2.TokenEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	if query != nil {
		if enconterError = model.CreateToken(query.SignedTo, query.Description); enconterError != nil {
			query2.APIResponse(c, enconterError, nil)
			return
		}
	} else {
		query2.SuccessResponse(c, query2.QUERY_NULL, nil)
	}
	query2.SuccessResponse(c, query2.OK, nil)
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
// @securityDefinitions.apikey BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /security/tokens [get]
func (t *Token) listToken(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	var list map[string]interface{}

	if list, enconterError = model.GetTokens(int(query.Offset), int(query.Limit), query.Sort); enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, list)
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
	query := &query2.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteTokenWithID(query.ID)
	if enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, nil)
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
	query := &query2.TokenEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = model.UpdateTokenWithID(query.ID, query.SignedTo, query.Description, query.IsUpdate); enconterError != nil {
			query2.API409Response(c, enconterError)
			return
		}

		query2.SuccessResponse(c, query2.OK, nil)
		return
	}
	query2.APIResponse(c, errors.New("invaild id"), nil)
}
