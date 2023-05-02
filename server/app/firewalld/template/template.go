package template

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/server/apis"
	templateModel "github.com/cylonchau/firewalld-gateway/utils/model"
)

type Template struct{}

func (t *Template) RegisterTemplateAPI(g *gin.RouterGroup) {
	g.GET("/", t.listTemplate)
	g.POST("/", t.createTemplate)
	g.DELETE("/", t.deleteTemplateWithID)
	g.PUT("/", t.updateTemplateWithID)
	g.GET("/port", t.listPort)
	g.POST("/port", t.createPort)
	g.DELETE("/port", t.deletePortWithID)
	g.PUT("/port", t.updatePortWithID)
	g.GET("/rich", t.listRich)
	g.POST("/rich", t.createRich)
	g.DELETE("/rich", t.deleteRichWithID)
	g.PUT("/rich", t.updateRichWithID)
}

func (t *Template) createTemplate(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.TemplateEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = templateModel.CreateTemplate(query.Name, query.Description, query.Target); enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	apis.SuccessResponse(c, apis.OK, nil)
}

func (t *Template) listTemplate(c *gin.Context) {
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
	if query.Simple == 0 {
		list, enconterError = templateModel.GetTemplates(int(query.Offset), int(query.Limit), query.Sort)
	} else {
		list, enconterError = templateModel.GetSimpleTemplates(int(query.Offset), 9999, query.Sort)
	}
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, list)
}

func (t *Template) deleteTemplateWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = templateModel.DeleteTemplateWithID(query.ID)
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, nil)
}

func (t *Template) updateTemplateWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.TemplateEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = templateModel.UpdateTemplateWithID(query.ID, query.Name, query.Description, query.Target); enconterError != nil {
			apis.API409Response(c, enconterError)
			return
		}

		apis.SuccessResponse(c, apis.OK, nil)
		return
	}
	apis.APIResponse(c, errors.New("invaild id"), nil)
}
