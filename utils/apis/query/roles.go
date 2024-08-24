package query

type RoleEditQuery struct {
	ID          int    `form:"id" json:"id,omitempty" binding:"omitempty"`
	Name        string `form:"name" json:"name,omitempty" binding:"required"`
	Description string `form:"description" json:"description"`
	RouterIDs   []int  `json:"router_ids" form:"router_ids"`
}
