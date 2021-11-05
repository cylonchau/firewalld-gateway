package query

import "firewall-api/utils/dbus"

type Query struct {
	Ip      string            `form:"ip" json:"ip" binding:"required"`
	Zone    string            `form:"zone,default=public" json:"zone"`
	Timeout int               `form:"timeout,default=0" json:"timeout"`
	Port    *dbus.Port        `form:"port" json:"port,omitempty"`
	Forward *dbus.ForwardPort `form:"forward" json:"forward,omitempty"`
	Rich    *dbus.Rule        `form:"rich" json:"rich,omitempty"`
	Service string            `form:"service" json:"service,omitempty"`
}

type PortQuery struct {
	Ip      string     `form:"ip" json:"ip" binding:"required"`
	Zone    string     `form:"zone,default=public" json:"zone"`
	Timeout int        `form:"timeout,default=0" json:"timeout"`
	Port    *dbus.Port `form:"port" json:"port,omitempty" binding:"required"`
}

type ForwardQuery struct {
	Ip      string            `form:"ip" json:"ip" binding:"required"`
	Zone    string            `form:"zone,default=public" json:"zone"`
	Timeout int               `form:"timeout,default=0" json:"timeout"`
	Forward *dbus.ForwardPort `form:"forward" json:"forward,omitempty" binding:"required"`
}

type ServiceQuery struct {
	Ip      string `form:"ip" json:"ip" binding:"required"`
	Zone    string `form:"zone,default=public" json:"zone"`
	Timeout int    `form:"timeout,default=0" json:"timeout"`
	Service string `form:"service" json:"service,omitempty" binding:"required"`
}

type RichQuery struct {
	Ip      string     `form:"ip" json:"ip" binding:"required"`
	Zone    string     `form:"zone,default=public" json:"zone"`
	Timeout int        `form:"timeout,default=0" json:"timeout"`
	Rich    *dbus.Rule `form:"rich" json:"rich,omitempty" binding:"required"`
}

type ServiceSettingQuery struct {
	Ip      string               `form:"ip" json:"ip" binding:"required"`
	Name    string               `form:"name" json:"name" binding:"required"`
	Setting *dbus.ServiceSetting `form:"setting" json:"setting,omitempty" binding:"required"`
}

type ZoneSettingQuery struct {
	Ip      string              `form:"ip" json:"ip" binding:"required"`
	Setting *dbus.QuerySettings `form:"setting" json:"Setting,omitempty" binding:"required"`
}

type RemoveQuery struct {
	Ip   string `form:"ip" json:"ip" binding:"required"`
	Name string `form:"name" json:"name" binding:"required"`
}
