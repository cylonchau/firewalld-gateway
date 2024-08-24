package tag

import (
	"errors"

	"github.com/gin-gonic/gin"

	query2 "github.com/cylonchau/firewalld-gateway/utils/apis/query"
	tagModel "github.com/cylonchau/firewalld-gateway/utils/model"
)

type Tag struct{}

func (t *Tag) RegisterTagAPI(g *gin.RouterGroup) {
	g.POST("/", t.createTag)
	g.GET("/", t.listTag)
	g.DELETE("/", t.deleteTagWithID)
	g.PUT("/", t.updateTagWithID)
}

func (t *Tag) createTag(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tagQuery := &query2.TagEditQuery{}
	enconterError = c.ShouldBindJSON(&tagQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = tagModel.CreateTag(tagQuery); enconterError != nil {
		query2.API409Response(c, enconterError)
		return
	}

	query2.SuccessResponse(c, query2.OK, nil)
}

func (t *Tag) listTag(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.ListQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := tagModel.GetTags(int(query.Offset), int(query.Limit), query.Sort)
	if enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, list)
}

func (t *Tag) deleteTagWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = tagModel.DeleteTagWithID(query.ID)
	if enconterError != nil {
		query2.API500Response(c, enconterError)
		return
	}
	query2.SuccessResponse(c, query2.OK, nil)
}

func (h *Tag) updateTagWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &query2.TagEditQuery{}
	enconterError = c.ShouldBindJSON(&query)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query2.APIResponse(c, enconterError, nil)
		return
	}

	if query.ID > 0 {
		if enconterError = tagModel.UpdateTagWithID(query); enconterError != nil {
			query2.API409Response(c, enconterError)
			return
		}

		query2.SuccessResponse(c, query2.OK, nil)
		return
	}
	query2.APIResponse(c, errors.New("invaild id"), nil)
}
