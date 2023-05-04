package apis

type UserQuery struct {
	Username string `form:"username" json:"username,omitempty" binding:"required"`
	Password string `form:"password" json:"password,omitempty" binding:"required"`
}

type UserResp struct {
	UserID   uint64 `form:"user_id" json:"user_id,omitempty" binding:"required"`
	Username string `form:"username" json:"username,omitempty" binding:"required"`
	Token    string `form:"token" json:"token,omitempty" binding:"required"`
	Name     string `form:"name" json:"name,omitempty" binding:"required"`
	LoginIP  string `form:"login_ip" json:"login_ip,omitempty"`
}

type InfoQuery struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type InfoResp struct {
	Username string `form:"username" json:"username" binding:"required"`
	UserRole string `form:"role" json:"role" binding:"required"`
}
