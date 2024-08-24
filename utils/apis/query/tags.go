package query

type TagEditQuery struct {
	Name        string `form:"name" json:"name,omitempty" binding:"required"`
	Description string `form:"description" json:"description"`
	ID          int    `form:"id" json:"id,omitempty" binding:"omitempty"`
}

type ListTagQuery struct {
	ID     int    `form:"id" json:"id"`
	Limit  uint16 `form:"limit,default=50" json:"limit"`
	Offset uint16 `form:"offset,default=0" json:"offset"`
	Sort   string `form:"sort,default=desc" json:"sort"`
}
