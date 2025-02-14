package firewalld

import (
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/api"
	"github.com/cylonchau/firewalld-gateway/config"
)

var (
	PORT          = "55556"
	InterfaceName = "com.github.cylonchau."
)

type DbusClientSerivce struct {
	client         *dbus.Conn
	defaultZone    string
	ip             string
	port           string
	eventLogFormat logFormat
}

func NewDbusClientService(addr string) (*DbusClientSerivce, error) {
	var (
		encounterError error
		conn           *dbus.Conn
		reply          dbus.RequestNameReply
	)
	if config.CONFIG.Port == "" {
		klog.V(5).Infof("Start connect to D-Bus service: %s:%s", addr, PORT)
	} else {
		PORT = config.CONFIG.DbusPort
		klog.V(5).Infof("Start connect to D-Bus service: %s:%s", addr, PORT)
	}

	if conn, encounterError = dbus.Connect("tcp:host="+addr+",port="+PORT, dbus.WithAuth(dbus.AuthAnonymous())); encounterError == nil {
		if encounterError == nil {
			appNameStr := strings.Split(config.CONFIG.AppName, " ")
			var registionName = InterfaceName + appNameStr[0]
			reply, encounterError = conn.RequestName(registionName, dbus.NameFlagDoNotQueue)
			switch reply {
			case dbus.RequestNameReplyInQueue:
				klog.Warningf("Interface %s already taken cannot be assigned again.", registionName)
			case dbus.RequestNameReplyExists:
				klog.Warningf("Interface %s cannot be assigned, because it's already taken by another owner", registionName)
			case dbus.RequestNameReplyAlreadyOwner:
				klog.Warningf("You are already the owner of %s. no need to ask again.", registionName)
			}
			if encounterError == nil {
				obj := conn.Object(api.INTERFACE, api.PATH)
				call := obj.Call(api.INTERFACE_GETDEFAULTZONE, dbus.FlagNoAutoStart)
				encounterError = call.Err
				if encounterError == nil {
					return &DbusClientSerivce{
						conn,
						call.Body[0].(string),
						addr,
						PORT,
						logFormat{},
					}, encounterError
				}
			}
		}
	}
	if encounterError != nil && conn != nil {
		conn.Close()
	}
	klog.Errorf("Connect to firewalld service failed: %v", encounterError)
	return nil, encounterError
}

/*
 * @title         Destroy
 * @description   off firewalld connection.
 * @middlewares          author    2021-10-31
 */
func (c *DbusClientSerivce) Destroy() {
	if c.client.Connected() {
		err := c.client.Close()
		if err != nil {
			klog.Errorf("Close D-Bus connection failed, %v", err)
		}
	}
}

/************************************************** fw service area ***********************************************************/

// @title         Reload
// @description   reload firewalld on runtime
// @middlewares   author           2024-08-07
// @return        error            error          "Possible errors:
//
//	ALREADY_ENABLED"
func (c *DbusClientSerivce) Reload() error {
	obj := c.client.Object(api.INTERFACE, api.PATH)
	c.printPath(api.INTERFACE_RELOAD)
	klog.V(4).Infof("Try to reload firewalld runtime.")
	call := obj.Call(api.INTERFACE_RELOAD, dbus.FlagNoAutoStart)

	if call.Err != nil {
		klog.Errorf("Reload firewalld failed: %v", call.Err.Error())
		return call.Err
	}
	klog.V(4).Infof("Reload firewalld success")
	return nil
}

/*
 * @title         flush currently zone zoneSettings to default zoneSettings.
 * @description   temporary Add rich language rule into zone.
 * @middlewares   author           2021-10-05
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RuntimeFlush(zone string) (encounterError error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	defaultZoneSetting := api.Settings{
		Target:      "accpet",
		Description: "reset by " + config.CONFIG.AppName,
		Short:       "public",
		Interface:   nil,
		Service: []string{
			"ssh",
			"dhcpv6-client",
		},
		Port: []*api.Port{
			{
				Port:     config.CONFIG.DbusPort,
				Protocol: "tcp",
			},
		},
	}

	var path dbus.ObjectPath
	if path, encounterError = c.generatePath(zone, api.ZONE_PATH); encounterError == nil {
		obj := c.client.Object(api.INTERFACE, path)
		c.printPath(api.CONFIG_UPDATE)
		klog.V(4).Infof("Try to flush current active zone (%s).", zone)
		call := obj.Call(api.CONFIG_UPDATE, dbus.FlagNoAutoStart, defaultZoneSetting)
		encounterError = call.Err
		if encounterError == nil || len(call.Body) <= 0 {
			if encounterError = c.Reload(); encounterError == nil {
				return nil
			}
		}
	}

	klog.Errorf("Flush current zone (%s) failed: %v", zone, encounterError)
	return encounterError
}

func (c *DbusClientSerivce) RuntimeSet(setting api.Settings) (encounterError error) {
	fmt.Println(setting)
	zone := c.GetDefaultZone()

	var path dbus.ObjectPath
	if path, encounterError = c.generatePath(zone, api.ZONE_PATH); encounterError == nil {
		obj := c.client.Object(api.INTERFACE, path)
		c.printPath(api.CONFIG_UPDATE)
		klog.V(4).Infof("Try to flush current active zone (%s).", zone)
		call := obj.Call(api.CONFIG_UPDATE, dbus.FlagNoAutoStart, setting)
		encounterError = call.Err
		if encounterError == nil || len(call.Body) <= 0 {
			if encounterError = c.Reload(); encounterError == nil {
				return nil
			}
		}
	}

	klog.Errorf("Set current zone (%s) failed: %v", zone, encounterError)
	return encounterError
}

// :title         Reload
// :description   genarate dbus path
// :Create        author   2021-10-05
// :Update        author   2024-09-06
// :return        error  error   "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) generatePath(zone, interfacePath string) (dbus.ObjectPath, error) {
	zoneid := c.getZoneId(zone)
	if zoneid < 0 {
		klog.Errorf("Invalid zone:", zone)
		return "", errors.New("Invalid zone " + interfacePath + "/" + zone)
	}
	p := fmt.Sprintf("%s/%d", interfacePath, zoneid)
	klog.V(5).Infof("D-Bus PATH: %s", p)
	return dbus.ObjectPath(p), nil
}
