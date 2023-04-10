package apis

type TagEditQuery struct {
	Name        string `form:"name" json:"name,omitempty" binding:"required"`
	Description string `form:"description" json:"description"`
}

type ListTagQuery struct {
	ID     int    `form:"id" json:"id"`
	Limit  uint16 `form:"limit,default=50" json:"limit"`
	Offset uint16 `form:"offset,default=0" json:"offset"`
}
