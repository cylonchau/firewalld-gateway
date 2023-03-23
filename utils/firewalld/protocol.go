package firewalld

import (
	"github.com/godbus/dbus/v5"

	"github.com/cylonchau/firewalld-gateway/apis"
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
func (c *DbusClientSerivce) AddProtocol(zone, protocol string, timeout uint32) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "protocol"
	c.eventLogFormat.resource = protocol
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_ADDPROTOCOL)
	call := obj.Call(apis.ZONE_ADDPROTOCOL, dbus.FlagNoAutoStart, zone, protocol, timeout)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError != nil {
		c.eventLogFormat.Format = CreateResourceFailedFormat
		c.printResourceEventLog()
		return call.Err
	}
	c.eventLogFormat.Format = CreateResourceSuccessFormat
	c.printResourceEventLog()
	return nil
}
