package template

import (
	"errors"

	"github.com/gin-gonic/gin"

	query2 "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	Model "github.com/cylonchau/firewalld-gateway/utils/model"
)

func (t *Template) createPort(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.PortEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = Model.CreatePort(query.Port, query.Protocol, query.TemplateId); enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	query2.SuccessResponse(c, query2.OK, nil)
}

func (t *Template) listPort(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := Model.GetPorts(int(query.Offset), int(query.Limit), query.Sort)
	if enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, list)
}

func (t *Template) deletePortWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = Model.DeletePortWithID(query.ID)
	if enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, nil)
}

func (t *Template) updatePortWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.PortEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = Model.UpdatePortWithID(query.ID, query.Port, query.Protocol, query.TemplateId); enconterError != nil {
			query2.API409Response(c, enconterError)
			return
		}

		query2.SuccessResponse(c, query2.OK, nil)
		return
	}
	query2.APIResponse(c, errors.New("invaild id"), nil)
}
