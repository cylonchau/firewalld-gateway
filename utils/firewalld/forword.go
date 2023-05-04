package firewalld

import (
	"fmt"
	"net"

	"github.com/godbus/dbus/v5"

	"github.com/cylonchau/firewalld-gateway/apis"
)

/************************************************** ForwardPort area ***********************************************************/

/*
 * @title         GetForwardPort
 * @description   temporary get IPv4 forward port in zone.
 * @middlewares      	  author           2021-10-27
 * @param         zone    		   string         	"If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        forwardPort      set          	"Return array of IPv4 forward ports previously added into zone.
 * @return        error            error          	"Possible errors:
 * 														INVALID_ZONE
 */
func (c *DbusClientSerivce) Listforwards(zone string) ([]apis.ForwardPort, error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	c.printPath(apis.ZONE_GETFORWARDPORT)
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	call := obj.Call(apis.ZONE_GETFORWARDPORT, dbus.FlagNoAutoStart, zone)

	c.eventLogFormat.encounterError = call.Err
	var forwards []apis.ForwardPort
	if c.eventLogFormat.encounterError == nil && len(call.Body) >= 0 {
		list, ok := call.Body[0].([][]string)
		if ok {
			for _, value := range list {
				var forward apis.ForwardPort
				forward, c.eventLogFormat.encounterError = apis.SliceToStruct(value)
				if c.eventLogFormat.encounterError == nil {
					forwards = append(forwards, forward)
				} else {
					break
				}
			}
		}
		c.eventLogFormat.Format = ListResourceSuccessFormat
		c.eventLogFormat.resource = len(forwards)
		c.printResourceEventLog()
		return forwards, nil
	}

	c.eventLogFormat.Format = ListResourceFailedFormat
	c.printResourceEventLog()
	return nil, c.eventLogFormat.encounterError

}

/*
 * @title         PermanentGetForwardPort
 * @description   permanent get IPv4 forward port in zone.
 * @middlewares      	  author           2021-10-29
 * @param         zone    		   string         	"If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        forwardPort      set          	"Return array of IPv4 forward ports previously added into zone.
 * @return        error            error          	"Possible errors:
 * 														INVALID_ZONE
 */
func (c *DbusClientSerivce) PermanentGetForwardPort(zone string) ([]apis.ForwardPort, error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = ListPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil

	var forwards []apis.ForwardPort
	var path dbus.ObjectPath

	if path, c.eventLogFormat.encounterError = c.generatePath(zone, apis.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)

		c.printResourceEventLog()

		c.printPath(apis.CONFIG_GETFORWARDPORT)
		call := obj.Call(apis.CONFIG_GETFORWARDPORT, dbus.FlagNoAutoStart)

		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil && len(call.Body) >= 0 {
			lists, ok := call.Body[0].([][]interface{})
			if ok {
				for _, value := range lists {
					var forward apis.ForwardPort
					if forward, c.eventLogFormat.encounterError = apis.SliceToStruct(value); c.eventLogFormat.encounterError == nil {
						forwards = append(forwards, forward)
					} else {
						break
					}
				}
				if c.eventLogFormat.encounterError == nil {
					c.eventLogFormat.Format = ListPermanentResourceSuccessFormat
					c.eventLogFormat.resource = len(forwards)
					c.printResourceEventLog()
					return forwards, nil
				}
			}
		}
	}

	c.eventLogFormat.Format = ListPermanentResourceFailedFormat
	c.printResourceEventLog()
	return nil, c.eventLogFormat.encounterError
}

/*
 * @title         AddForwardPort
 * @description   temporary Add the IPv4 forward port into zone.
 * @middlewares      	  author           2021-09-29
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
 * @return        error            error          "Possible errors:
 * 													INVALID_ZONE,
 * 													INVALID_PORT,
 * 													MISSING_PROTOCOL,
 * 													INVALID_PROTOCOL,
 * 													INVALID_ADDR,
 * 													INVALID_FORWARD,
 * 													ALREADY_ENABLED,
 * 													INVALID_COMMAND"
 */
func (c *DbusClientSerivce) AddForwardPort(zone string, timeout uint32, forward *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = fmt.Sprintf("%s => %s:%s", forward.Port, forward.ToAddr, forward.ToPort)

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printResourceEventLog()

	c.printPath(apis.ZONE_ADDFORWARDPORT)
	call := obj.Call(apis.ZONE_ADDFORWARDPORT, dbus.FlagNoAutoStart, zone, forward.Port, forward.Protocol, forward.ToPort, forward.ToAddr, timeout)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil && len(call.Body) > 0 {
		c.eventLogFormat.Format = CreateResourceSuccessFormat
		c.printResourceEventLog()
		return nil
	}
	c.eventLogFormat.Format = CreatePermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

/*
 * @title         PermanentAddForwardPort
 * @description   temporary Add the IPv4 forward port into zone.
 * @middlewares      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) AddPermanentForwardPort(zone string, forward *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = CreatePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = fmt.Sprintf("%s => %s:%s", forward.Port, forward.ToAddr, forward.ToPort)

	path, _ := c.generatePath(zone, apis.ZONE_PATH)
	obj := c.client.Object(apis.INTERFACE, path)
	c.printResourceEventLog()

	c.printPath(apis.CONFIG_ZONE_ADDFORWARDPORT)
	call := obj.Call(apis.CONFIG_ZONE_ADDFORWARDPORT, dbus.FlagNoAutoStart, forward.Port, forward.Protocol, forward.ToPort, forward.ToAddr)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil {
		c.eventLogFormat.Format = CreatePermanentResourceSuccessFormat
		c.printResourceEventLog()
		return nil
	}

	c.eventLogFormat.Format = CreatePermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

/*
 * @title         RemoveForwardPort
 * @description   temporary (runtime) Remove IPv4 forward port ((port, protocol, toport, toaddr)) from zone.
 * @middlewares      	  author           2021-09-29
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													INVALID_ZONE,
 * 													INVALID_PORT,
 * 													MISSING_PROTOCOL,
 * 													INVALID_PROTOCOL,
 * 													INVALID_ADDR,
 * 													INVALID_FORWARD,
 * 													ALREADY_ENABLED,
 * 													INVALID_COMMAND"
 */
func (c *DbusClientSerivce) RemoveForwardPort(zone string, forward *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = RemoveResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = fmt.Sprintf("%s => %s:%s", forward.Port, forward.ToAddr, forward.ToPort)

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	c.printResourceEventLog()

	c.printPath(apis.ZONE_REMOVEFORWARDPORT)
	call := obj.Call(apis.ZONE_REMOVEFORWARDPORT, dbus.FlagNoAutoStart, zone, forward.Port, forward.Protocol, forward.ToPort, forward.ToAddr)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil {
		c.eventLogFormat.Format = RemoveResourceSuccessFormat
		c.printResourceEventLog()
		return nil
	}
	c.eventLogFormat.Format = RemoveResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

/*
 * @title         PermanentRemoveForwardPort
 * @description   Permanently remove (port, protocol, toport, toaddr) from list of forward ports of zone.
 * @middlewares      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RemovePermanentForwardPort(zone string, forward *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = RemovePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = fmt.Sprintf("%s => %s:%s", forward.Port, forward.ToAddr, forward.ToPort)

	var path dbus.ObjectPath
	path, c.eventLogFormat.encounterError = c.generatePath(zone, apis.ZONE_PATH)

	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)

		c.printResourceEventLog()
		c.printPath(apis.CONFIG_ZONE_REMOVEFORWARDPORT)
		call := obj.Call(apis.CONFIG_ZONE_REMOVEFORWARDPORT, dbus.FlagNoAutoStart, forward.Port, forward.Protocol, forward.ToPort, forward.ToAddr)

		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = RemovePermanentResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}

	c.eventLogFormat.Format = RemovePermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

/*
 * @title         QueryForwardPort
 * @description   temporary (runtime) query whether the IPv4 forward port (port, protocol, toport, toaddr) has been added into zone.
 * @middlewares      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													INVALID_ZONE,
 * 													INVALID_PORT,
 * 													MISSING_PROTOCOL,
 * 													INVALID_PROTOCOL,
 * 													INVALID_ADDR,
 * 													INVALID_FORWARD,
 * 													ALREADY_ENABLED,
 * 													INVALID_COMMAND"
 */
func (c *DbusClientSerivce) QueryForwardPort(zone, portProtocol, toHostPort string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = fmt.Sprintf("%s/%s", portProtocol, toHostPort)

	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, enconterError := net.SplitHostPort(toHostPort)
	c.eventLogFormat.encounterError = enconterError

	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, apis.PATH)
		c.printResourceEventLog()

		c.printPath(apis.ZONE_QUERYFORWARDPORT)
		call := obj.Call(apis.ZONE_QUERYFORWARDPORT, dbus.FlagNoAutoStart, zone, port, protocol, toPort, toAddr)
		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil || call.Body[0].(bool) {
			c.eventLogFormat.Format = QueryResourceSuccessFormat
			c.printResourceEventLog()
			return true
		}
	}
	c.eventLogFormat.Format = QueryResourceFailedFormat
	c.printResourceEventLog()
	return false
}

/*
 * @title         PermanentQueryForwardPort
 * @description   Permanently remove (port, protocol, toport, toaddr) from list of forward ports of zone.
 * @middlewares      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentQueryForwardPort(zone, portProtocol, toHostPort string) (b bool) {
	var enconterError error
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = QueryPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "NAT forward"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = fmt.Sprintf("%s/%s", portProtocol, toHostPort)

	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, enconterError := net.SplitHostPort(toHostPort)
	c.eventLogFormat.encounterError = enconterError

	if c.eventLogFormat.encounterError == nil {
		var path dbus.ObjectPath
		if path, c.eventLogFormat.encounterError = c.generatePath(zone, apis.ZONE_PATH); c.eventLogFormat.encounterError == nil {
			obj := c.client.Object(apis.INTERFACE, path)

			c.printResourceEventLog()
			c.printPath(apis.CONFIG_ZONE_QUERYFORWARDPORT)
			call := obj.Call(apis.CONFIG_ZONE_QUERYFORWARDPORT, dbus.FlagNoAutoStart, port, protocol, toPort, toAddr)
			c.eventLogFormat.encounterError = call.Err

			if enconterError == nil || call.Body[0].(bool) {
				c.eventLogFormat.Format = QueryPermanentResourceSuccessFormat
				c.printResourceEventLog()
				return true
			}
		}
	}
	c.eventLogFormat.Format = QueryPermanentResourceFailedFormat
	c.printResourceEventLog()
	return false
}
