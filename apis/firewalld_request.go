package apis

type Query struct {
	Ip      string       `form:"ip" json:"ip" binding:"required"`
	Zone    string       `form:"zone,default=public" json:"zone"`
	Timeout int          `form:"timeout,default=0" json:"timeout"`
	Port    *Port        `form:"port" json:"port,omitempty"`
	Forward *ForwardPort `form:"forward" json:"forward,omitempty"`
	Rich    *Rule        `form:"rich" json:"rich,omitempty"`
	Service string       `form:"service" json:"service,omitempty"`
}

type PortQuery struct {
	Ip      string `form:"ip" json:"ip" binding:"required"`
	Zone    string `form:"zone,default=public" json:"zone"`
	Timeout int    `form:"timeout,default=0" json:"timeout"`
	Port    *Port  `form:"port" json:"port,omitempty" binding:"required"`
}

type ForwardQuery struct {
	Ip      string       `form:"ip" json:"ip" binding:"required"`
	Zone    string       `form:"zone,default=public" json:"zone"`
	Timeout int          `form:"timeout,default=0" json:"timeout"`
	Forward *ForwardPort `form:"forward" json:"forward,omitempty" binding:"required"`
}

type ServiceQuery struct {
	Ip      string `form:"ip" json:"ip" binding:"required"`
	Zone    string `form:"zone,default=public" json:"zone"`
	Timeout int    `form:"timeout,default=0" json:"timeout"`
	Service string `form:"service" json:"service,omitempty" binding:"required"`
}

type RichQuery struct {
	Ip      string `form:"ip" json:"ip" binding:"required"`
	Zone    string `form:"zone,default=public" json:"zone"`
	Timeout int    `form:"timeout,default=0" json:"timeout"`
	Rich    *Rule  `form:"rich" json:"rich,omitempty" binding:"required"`
}

type ServiceSettingQuery struct {
	Ip      string          `form:"ip" json:"ip" binding:"required"`
	Name    string          `form:"name" json:"name" binding:"required"`
	Setting *ServiceSetting `form:"setting" json:"setting,omitempty" binding:"required"`
}

type ZoneSettingQuery struct {
	Ip      string         `form:"ip" json:"ip" binding:"required"`
	Setting *QuerySettings `form:"setting" json:"Setting,omitempty" binding:"required"`
}

type RemoveQuery struct {
	Ip   string `form:"ip" json:"ip" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
}
