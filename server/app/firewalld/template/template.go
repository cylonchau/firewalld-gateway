package template

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/praserx/ipconv"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

type Template struct{}

func (t *Template) RegisterTemplateAPI(g *gin.RouterGroup) {
	g.GET("/", t.listTemplate)
	g.POST("/:id", t.getTemplateRulesWithID)
	g.PUT("/", t.createTemplate)
	g.DELETE("/", t.deleteTemplateWithID)
	g.POST("/", t.updateTemplateWithID)
	g.GET("/port", t.listPort)
	g.POST("/port", t.createPort)
	g.DELETE("/port", t.deletePortWithID)
	g.PUT("/port", t.updatePortWithID)
	g.GET("/rich", t.listRich)
	g.PUT("/rich", t.createRich)
	g.DELETE("/rich", t.deleteRichWithID)
	g.POST("/rich", t.updateRichWithID)
}

// createTemplate godoc
// @Summary Create a new firewalld template.
// @Description Create a new firewalld template.
// @Tags Template
// @Accept  json
// @Produce json
// @Param query body query.TemplateEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template [PUT]
func (t *Template) createTemplate(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templateQuery := &query.TemplateEditQuery{}
	enconterError = c.ShouldBindJSON(&templateQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = model.CreateTemplate(templateQuery.Name, templateQuery.Description, templateQuery.Target); enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	query.SuccessResponse(c, query.OK, nil)
}

// listTemplate godoc
// @Summary List template.
// @Description List template.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.ListQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template [GET]
func (t *Template) listTemplate(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templateQuery := &query.ListQuery{}
	enconterError = c.Bind(&templateQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	var list map[string]interface{}
	if templateQuery.Simple == 0 {
		list, enconterError = model.GetTemplates(templateQuery.Title, int(templateQuery.Offset), int(templateQuery.Limit), templateQuery.Sort)
	} else {
		list, enconterError = model.GetSimpleTemplates(int(templateQuery.Offset), 9999, templateQuery.Sort)
	}
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, list)
}

// getTemplateRulesWithID godoc
// @Summary List template rules.
// @Description List template rules by template ID.
// @Tags Template
// @Accept json
// @Produce json
// @Param template_id path int true "Template ID"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/{id} [POST]
func (t *Template) getTemplateRulesWithID(c *gin.Context) {
	// 1. 获取参数
	templateQuery := &query.QueryWithID{}
	var enconterError error
	if enconterError = c.ShouldBindUri(templateQuery); enconterError != nil {
		query.API400Response(c, enconterError)
		return
	}

	templateDetails, enconterError := model.GetRichWithDetailsByTemplateID(templateQuery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	if hosts, enconterError := model.GetHostsByTagName(templateDetails.Short); enconterError == nil {
		if len(hosts) == 0 {
			query.API404Response(c, fmt.Errorf("No host in tag %s", templateDetails.Short))
			return
		}
		for _, host := range hosts {
			ip := ipconv.IntToIPv4(host.IP).String()

			dbusClient, enconterError := firewalld.NewDbusClientService(ip)
			if enconterError != nil {
				query.ConnectDbusService(c, enconterError)
				return
			}
			defer dbusClient.Destroy()
			if err := dbusClient.RuntimeSet(*templateDetails); err != nil {
				query.API500Response(c, err)
				return
			}
		}
		query.SuccessResponse(c, query.OK, templateDetails)
		return
	}
	// 3. 返回成功响应
	query.SuccessResponse(c, query.OK, templateDetails)
}

// deleteTemplateWithID godoc
// @Summary Delete template with template id.
// @Description Delete template with template id.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.IDQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template [DELETE]
func (t *Template) deleteTemplateWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templateQuery := &query.IDQuery{}
	enconterError = c.Bind(&templateQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteTemplateWithID(templateQuery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// updateTemplateWithID godoc
// @Summary Update template information with template id.
// @Description Update template information with template id.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.TemplateEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template [POST]
func (t *Template) updateTemplateWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templateQuery := &query.TemplateEditQuery{}
	enconterError = c.ShouldBindJSON(&templateQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if templateQuery.ID > 0 {
		if enconterError = model.UpdateTemplateWithID(templateQuery.ID, templateQuery.Name, templateQuery.Description, templateQuery.Target); enconterError != nil {
			query.API409Response(c, enconterError)
			return
		}

		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}
