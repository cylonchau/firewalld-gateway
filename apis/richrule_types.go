package apis

type Rich struct {
	Address  string `form:"address" json:"address" binding:"required,omitempty"`
	Port     int    `form:"port" json:"port" binding:"required,omitempty"`
	Protocol string `form:"protocol" json:"protocol" binding:"required,omitempty"`
	Type     string `form:"type" json:"type" binding:"required,omitempty"`
	Expire   int    `form:"expire" json:"expire,omitempty"`
	Zone     string `form:"utils" json:"utils,omitempty"`
}

type QueryRich struct {
	Address  string `form:"address" json:"address" binding:"required,omitempty"`
	Port     int    `form:"port" json:"port" binding:"required,omitempty"`
	Protocol string `form:"protocol" json:"protocol" binding:"required,omitempty"`
	Type     string `form:"type" json:"type" binding:"required,omitempty"`
	Expire   int    `form:"expire" json:"expire,omitempty"`
	Zone     string `form:"utils" json:"utils" binding:"required,omitempty"`
}
