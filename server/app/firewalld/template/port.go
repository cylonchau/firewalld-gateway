package template

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

// createPort godoc
// @Summary Create a new port rule to template.
// @Description Create a new port rule to template.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.PortEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/port [PUT]
func (t *Template) createPort(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templatePortquery := &query.PortEditQuery{}
	enconterError = c.ShouldBindJSON(&templatePortquery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = model.CreatePort(templatePortquery.Port, templatePortquery.Protocol, templatePortquery.TemplateId); enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	query.SuccessResponse(c, query.OK, nil)
}

// listPort godoc
// @Summary List port rules.
// @Description List port rules.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.ListQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/port [GET]
func (t *Template) listPort(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templatePortquery := &query.ListQuery{}
	enconterError = c.Bind(&templatePortquery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := model.GetPorts(templatePortquery.Title, int(templatePortquery.Offset), int(templatePortquery.Limit), templatePortquery.Sort)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, list)
}

// deletePortWithID godoc
// @Summary Delete port rule with rule id.
// @Description Delete port rule with rule id.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.IDQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/port [DELETE]
func (t *Template) deletePortWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	templatePortquery := &query.IDQuery{}
	enconterError = c.Bind(&templatePortquery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeletePortWithID(templatePortquery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// updatePortWithID godoc
// @Summary Update port rule information with port rule id.
// @Description Update port rule information with port rule id.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.PortEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/port [POST]
func (t *Template) updatePortWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templatePortquery := &query.PortEditQuery{}
	enconterError = c.ShouldBindJSON(&templatePortquery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if templatePortquery.ID > 0 {
		if enconterError = model.UpdatePortWithID(templatePortquery.ID, templatePortquery.Port, templatePortquery.Protocol, templatePortquery.TemplateId); enconterError != nil {
			query.API409Response(c, enconterError)
			return
		}

		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}
