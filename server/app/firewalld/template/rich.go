package template

import (
	"errors"

	"github.com/gin-gonic/gin"

	query2 "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	Model "github.com/cylonchau/firewalld-gateway/utils/model"
)

func (t *Template) createRich(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.RichEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = Model.CreateRich(query); enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	query2.SuccessResponse(c, query2.OK, nil)
}

func (t *Template) listRich(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := Model.GetRich(int(query.Offset), int(query.Limit), query.Sort)
	if enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, list)
}

func (t *Template) deleteRichWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = Model.DeleteRichWithID(query.ID)
	if enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, nil)
}

func (t *Template) updateRichWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.RichEditQuery{}

	// 手动对请求参数进行详细的业务规则校验
	if enconterError = c.ShouldBindJSON(&query); enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = Model.UpdateRichWithID(query); enconterError != nil {
			query2.APIResponse(c, enconterError, nil)
			return
		}
		query2.SuccessResponse(c, query2.OK, nil)
		return
	}
	query2.APIResponse(c, errors.New("invaild id"), nil)
}
