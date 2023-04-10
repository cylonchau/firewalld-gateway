package apis

type HostQuery struct {
	IP     string `form:"ip" json:"ip,omitempty" binding:"required"`
	TagId  int    `form:"tag_id" json:"tag_id"  binding:"required"`
	Limit  uint16 `form:"limit,default=100" json:"limit"`
	Offset uint16 `form:"offset,default=0" json:"offset"`
}

type AsyncHostQuery struct {
	IPRange string `form:"ip_range" json:"ip_range,omitempty" binding:"required"`
	TagId   int    `form:"tag_id" json:"tag_id"  binding:"required"`
}

type ListHostQuery struct {
	Limit  uint16 `form:"limit,default=100" json:"limit"`
	Offset uint16 `form:"offset,default=0" json:"offset"`
}

type IDQuery struct {
	ID uint64 `form:"id" json:"id,omitempty" binding:"required"`
}
