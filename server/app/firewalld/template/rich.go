package template

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	"github.com/cylonchau/firewalld-gateway/utils/model"
)

// createRich godoc
// @Summary Create a new rich rule to template.
// @Description Create a new rich rule to template.
// @Tags Template
// @Accept  json
// @Produce json
// @Param query body query.RichEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/rich [PUT]
func (t *Template) createRich(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templateRichQuery := &query.RichEditQuery{}
	enconterError = c.ShouldBindJSON(&templateRichQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = model.CreateRich(templateRichQuery); enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	query.SuccessResponse(c, query.OK, nil)
}

// listRich godoc
// @Summary List rich rules.
// @Description List rich rules.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.ListQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/rich [GET]
func (t *Template) listRich(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templateRichQuery := &query.ListQuery{}
	enconterError = c.Bind(&templateRichQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := model.GetRich(templateRichQuery.Title, int(templateRichQuery.Offset), int(templateRichQuery.Limit), templateRichQuery.Sort)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, list)
}

// deleteRichWithID godoc
// @Summary Delete rich rule with rule id.
// @Description Delete rich rule with rule id.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.IDQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/rich [DELETE]
func (t *Template) deleteRichWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	templateRichQuery := &query.IDQuery{}
	enconterError = c.Bind(&templateRichQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = model.DeleteRichWithID(templateRichQuery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// updateRichWithID godoc
// @Summary Update rich rule information with rich rule id.
// @Description Update rich rule information with rich rule id.
// @Tags Template
// @Accept json
// @Produce json
// @Param query body query.RichEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/template/rich [POST]
func (t *Template) updateRichWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	templateRichQuery := &query.RichEditQuery{}

	// 手动对请求参数进行详细的业务规则校验
	if enconterError = c.ShouldBindJSON(&templateRichQuery); enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	if templateRichQuery.ID > 0 {
		if enconterError = model.UpdateRichWithID(templateRichQuery); enconterError != nil {
			query.APIResponse(c, enconterError, nil)
			return
		}
		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}
