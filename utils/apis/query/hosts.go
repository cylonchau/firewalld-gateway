package query

type HostQuery struct {
	IP       string `form:"ip" json:"ip,omitempty" binding:"required"`
	TagId    int    `form:"tag_id" json:"tag_id"  binding:"required"`
	Hostname string `form:"hostname" json:"hostname" `
	ID       int    `form:"id" json:"id" binding:"omitempty"`
	Limit    uint16 `form:"limit,default=100" json:"limit"`
	Offset   uint16 `form:"offset,default=0" json:"offset"`
}
type AsyncHostQuery struct {
	IPRange string `form:"ip_range" json:"ip_range,omitempty" binding:"required"`
	TagId   int    `form:"tag_id" json:"tag_id"  binding:"required"`
}

type ListHostQuery struct {
	Limit  uint16 `form:"limit,default=100" json:"limit"`
	Offset uint16 `form:"offset,default=0" json:"offset"`
	Sort   string `form:"sort,default=desc" json:"sort"`
}

type IDQuery struct {
	ID uint64 `form:"id" json:"id,omitempty" binding:"required"`
}
