package tag

import (
	"fmt"

	"github.com/gin-gonic/gin"

	tagModel "github.com/cylonchau/firewalld-gateway/model"
	"github.com/cylonchau/firewalld-gateway/server/apis"
)

type Tag struct{}

func (t *Tag) RegisterTagAPI(g *gin.RouterGroup) {
	g.POST("/", t.createTag)
	g.GET("/", t.listTag)
	g.DELETE("/", t.deleteTagWithID)
}

func (t *Tag) createTag(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tagQuery := &apis.TagEditQuery{}
	enconterError = c.ShouldBindJSON(&tagQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = tagModel.CreateTag(tagQuery); enconterError != nil {
		apis.API409Response(c, enconterError)
		return
	}

	apis.SuccessResponse(c, apis.OK, nil)
}

func (t *Tag) listTag(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.ListTagQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := tagModel.GetTags(int(query.Offset), int(query.Limit))
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, list)
}

func (t *Tag) deleteTagWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	query := &apis.IDQuery{}
	enconterError = c.Bind(&query)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		apis.APIResponse(c, enconterError, nil)
		return
	}
	fmt.Println(query)
	enconterError = tagModel.DeleteTagWithID(query.ID)
	if enconterError != nil {
		apis.API500Response(c, enconterError)
		return
	}
	apis.SuccessResponse(c, apis.OK, nil)
}
