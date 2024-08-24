package query

type UserQuery struct {
	Username string `form:"username" json:"username,omitempty" binding:"required"`
	Password string `form:"password" json:"password,omitempty" binding:"required"`
}

type AllocateRoleQuery struct {
	UserID  uint64 `form:"user_id" json:"user_id,omitempty" binding:"required"`
	RoleIDs []int  `json:"role_ids" form:"role_ids" binding:"required"`
}

type AllocateRouterQuery struct {
	UserID  uint64 `form:"user_id" json:"user_id,omitempty" binding:"required"`
	RoleIDs []int  `json:"router_ids" form:"router_ids" binding:"required"`
}

type UserResp struct {
	UserID       uint64   `form:"user_id" json:"user_id,omitempty" binding:"required"`
	Username     string   `form:"username" json:"username,omitempty" binding:"required"`
	Token        string   `form:"token" json:"token,omitempty" binding:"required"`
	Name         string   `form:"name" json:"name,omitempty" binding:"required"`
	IsPrivileged bool     `form:"is_privileged" json:"is_privileged,omitempty" binding:"required"`
	Roles        []string `form:"roles" json:"roles,omitempty" binding:"required"`
	LoginIP      string   `form:"login_ip" json:"login_ip,omitempty"`
}

type UserEditQuery struct {
	ID       int    `form:"id" json:"id,omitempty" binding:"omitempty"`
	Username string `form:"username" json:"username,omitempty" binding:"required"`
	Password string `form:"password" json:"password,omitempty"`
	RoleIDs  []int  `json:"role_ids" form:"role_ids"`
}

type InfoQuery struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type InfoResp struct {
	Username string `form:"username" json:"username" binding:"required"`
	UserRole string `form:"role" json:"role" binding:"required"`
}
