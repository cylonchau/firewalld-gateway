package firewalld

import (
	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
)

/************************************************** Protocol area ***********************************************************/

// @title         AddProtocol
// @description   temporary get a firewalld port list
// @auth      	  author           2021-09-29
// @param         zone    		   string         "e.g. public|dmz.. If zone is empty string, use default zone. "
// @param         protocol         string         "e.g. tcp|udp... The protocol can be any protocol supported by the system."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the protocol was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_PROTOCOL, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddProtocol(zone, protocol string, timeout int) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	printPath(apis.PATH, apis.ZONE_ADDPROTOCOL)

	klog.V(4).Infof("Trying to create protocol rule in zone %s/%s", zone, protocol)
	call := obj.Call(apis.ZONE_ADDPROTOCOL, dbus.FlagNoAutoStart, zone, protocol, timeout)

	if call.Err != nil {
		klog.Errorf("Create protocol rule failed: %v", call.Err.Error())
		return call.Err
	}
	return nil
}
