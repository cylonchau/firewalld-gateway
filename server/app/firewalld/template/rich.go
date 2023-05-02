package template

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/apis"
	Model "github.com/cylonchau/firewalld-gateway/utils/model"
)

func (t *Template) createRich(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.RichEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = Model.CreateRich(query); enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	apis.SuccessResponse(c, apis.OK, nil)
}

func (t *Template) listRich(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := Model.GetRich(int(query.Offset), int(query.Limit), query.Sort)
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, list)
}

func (t *Template) deleteRichWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = Model.DeleteRichWithID(query.ID)
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, nil)
}

func (t *Template) updateRichWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.RichEditQuery{}

	// 手动对请求参数进行详细的业务规则校验
	if enconterError = c.ShouldBindJSON(&query); enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = Model.UpdateRichWithID(query); enconterError != nil {
			apis.APIResponse(c, enconterError, nil)
			return
		}
		apis.SuccessResponse(c, apis.OK, nil)
		return
	}
	apis.APIResponse(c, errors.New("invaild id"), nil)
}
