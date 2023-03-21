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
	Port    Port   `form:"port" json:"port,omitempty" binding:"required"`
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
	Host        string          `form:"host" json:"host" binding:"required"`
	ServiceName string          `form:"service_name" json:"service_name" binding:"required"`
	Setting     *ServiceSetting `form:"setting" json:"setting,omitempty" binding:"required"`
}

type ZoneSettingQuery struct {
	Ip      string         `form:"ip" json:"ip" binding:"required"`
	Setting *QuerySettings `form:"setting" json:"Setting,omitempty" binding:"required"`
}

type RemoveQuery struct {
	Ip   string `form:"ip" json:"ip" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
}

type BatchPortQuery struct {
	Ports []PortQuery `form:"ports" json:"ports"`
}

type BatchSettingQuery struct {
	Hosts []string `form:"hosts" json:"hosts,omitempty" binding:"required"`
}

type ZoneDst struct {
	Host string `form:"host" json:"host,omitempty" binding:"required"`
	Zone string `form:"zone" json:"zone,omitempty" binding:"required"`
}

type BatchZoneQuery struct {
	ActionObject []ZoneDst `form:"action_object" json:"action_object,omitempty" binding:"required"`
}

type BatchServiceQuery struct {
	Services []ServiceQuery `form:"services" json:"services,omitempty"`
}
