package query

type ListQuery struct {
	ID     int    `form:"id" json:"id"`
	Simple uint8  `form:"simple,omitempty,default=0" json:"simple"`
	Limit  uint16 `form:"limit,default=10" json:"limit"`
	Offset uint16 `form:"offset,default=0" json:"offset"`
	Sort   string `form:"sort,default=desc" json:"sort"`
}

type LsQuery struct {
	IDs    []int  `form:"ids" json:"ids" binding:"required"`
	Simple uint8  `form:"simple,omitempty,default=0" json:"simple"`
	Limit  uint16 `form:"limit,default=10" json:"limit"`
	Offset uint16 `form:"offset,default=0" json:"offset"`
	Sort   string `form:"sort,default=desc" json:"sort"`
}
