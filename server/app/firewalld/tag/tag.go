package tag

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
	tagModel "github.com/cylonchau/firewalld-gateway/utils/model"
)

type Tag struct{}

func (t *Tag) RegisterTagAPI(g *gin.RouterGroup) {
	g.PUT("/", t.createTag)
	g.GET("/", t.listTag)
	g.DELETE("/", t.deleteTagWithID)
	g.POST("/", t.updateTagWithID)
}

// createTag godoc
// @Summary Add a new tag into uranus.
// @Description Add a new tag into uranus.
// @Tags Tags
// @Accept  json
// @Produce json
// @Param query body query.TagEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/tag [PUT]
func (t *Tag) createTag(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tagQuery := &query.TagEditQuery{}
	enconterError = c.ShouldBindJSON(&tagQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if enconterError = tagModel.CreateTag(tagQuery); enconterError != nil {
		query.API409Response(c, enconterError)
		return
	}

	query.SuccessResponse(c, query.OK, nil)
}

// listTag godoc
// @Summary Get tag list from uranus.
// @Description Get tag list from uranus.
// @Tags Tags
// @Accept  json
// @Produce json
// @Param query body query.ListQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/tag [GET]
func (t *Tag) listTag(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tagQuery := &query.ListQuery{}
	enconterError = c.Bind(&tagQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	list, enconterError := tagModel.GetTags(tagQuery.Title, int(tagQuery.Offset), int(tagQuery.Limit), tagQuery.Sort)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, list)
}

// deleteTagWithID godoc
// @Summary Delete tag with tag id.
// @Description Delete tag with tag id.
// @Tags Tags
// @Accept json
// @Produce json
// @Param query body query.IDQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/tag [DELETE]
func (t *Tag) deleteTagWithID(c *gin.Context) {

	// 1. 获取参数和参数校验
	var enconterError error
	tagQuery := &query.IDQuery{}
	enconterError = c.Bind(&tagQuery)
	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}
	enconterError = tagModel.DeleteTagWithID(tagQuery.ID)
	if enconterError != nil {
		query.API500Response(c, enconterError)
		return
	}
	query.SuccessResponse(c, query.OK, nil)
}

// updateTagWithID godoc
// @Summary Update tag information with tag id.
// @Description Update tag information with tag id.
// @Tags Tags
// @Accept  json
// @Produce json
// @Param query body query.TagEditQuery false "body"
// @Security BearerAuth
// @Success 200 {object} interface{}
// @Router /fw/tag [POST]
func (h *Tag) updateTagWithID(c *gin.Context) {
	// 1. 获取参数和参数校验
	var enconterError error
	tagQuery := &query.TagEditQuery{}
	enconterError = c.ShouldBindJSON(&tagQuery)

	// 手动对请求参数进行详细的业务规则校验
	if enconterError != nil {
		query.APIResponse(c, enconterError, nil)
		return
	}

	if tagQuery.ID > 0 {
		if enconterError = tagModel.UpdateTagWithID(tagQuery); enconterError != nil {
			query.API409Response(c, enconterError)
			return
		}

		query.SuccessResponse(c, query.OK, nil)
		return
	}
	query.APIResponse(c, errors.New("invaild id"), nil)
}
