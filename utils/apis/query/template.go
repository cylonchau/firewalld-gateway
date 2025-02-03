package query

type TemplateEditQuery struct {
	Name        string `form:"name" json:"name,omitempty" binding:"required"`
	Description string `form:"description" json:"description"`
	Target      string `form:"target" json:"target"`
	ID          uint64 `form:"id" json:"id,omitempty" binding:"omitempty"`
}

type PortEditQuery struct {
	Port       uint16 `form:"port" json:"port" binding:"required"`
	Protocol   string `form:"protocol" json:"protocol" binding:"required"`
	TemplateId int    `form:"template_id" json:"template_id" binding:"required"`
	ID         uint64 `form:"id" json:"id,omitempty" binding:"omitempty"`
}

type RichEditQuery struct {
	Family      string `form:"family" json:"family,omitempty" binding:"omitempty"`
	Source      string `form:"source" json:"source,omitempty" binding:"omitempty"`
	Destination string `form:"destination" json:"destination,omitempty" binding:"omitempty"`
	Port        string `form:"port" json:"port,omitempty" binding:"omitempty"`
	Protocol    string `form:"protocol" json:"protocol,omitempty" binding:"omitempty"`
	Action      string `form:"action" json:"action,omitempty" binding:"required"`
	Limit       uint16 `form:"limit" json:"limit,omitempty" binding:"omitempty"`
	LimitUnit   string `form:"limit_unit" json:"limit_unit,omitempty" binding:"omitempty"`
	TemplateId  int    `form:"template_id" json:"template_id" binding:"required"`
	ID          uint64 `form:"id" json:"id,omitempty" binding:"omitempty"`
}
