package apis

type TokenEditQuery struct {
	Description string `json:"description,omitempty" binding:"omitempty"`
	SignedTo    string `json:"signed_to,omitempty" binding:"omitempty"`
	IsUpdate    bool   `json:"is_update,omitempty" binding:"omitempty"`
	ID          uint64 `form:"id" json:"id,omitempty" binding:"omitempty"`
}
