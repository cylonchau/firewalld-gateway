package firewalld

import (
	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
)

/************************************************** port area ***********************************************************/

// @title         addPort
// @description   temporary add a firewalld port
// @auth      	  author           2021-09-29
// @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
// @param         zone    		   string         "e.g. public|dmz.. The empty string is usage default zone, is currently firewalld defualt zone"
// @param         timeout    	   int	          "Timeout, 0 is the permanent effect of the currently service startup state."
// @return        zoneName         string         "Returns name of zone to which the protocol was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_PORT, MISSING_PROTOCOL, INVALID_PROTOCOL, ALREADY_ENABLED, INVALID_COMMAND."
func (c *DbusClientSerivce) AddPort(port *apis.Port, zone string, timeout int) error {

	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_ADDPORT)
	klog.V(4).Infof("Trying create port rule in zone <%s>, %s/%s, lifecycle is %d", zone, port.Port, port.Protocol, timeout)
	call := obj.Call(apis.ZONE_ADDPORT, dbus.FlagNoAutoStart, zone, port.Port, port.Protocol, timeout)
	if call.Err != nil || len(call.Body) <= 0 {
		klog.Errorf("Create a Port Rule Failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

// @title         PermanentAddPort
// @description   Permanently add port & procotol to list of ports of zone.
// @auth      	  author           2021-09-29
// @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
// @param         zone    		   string         "e.g. public|dmz.. The empty string is usage default zone, is currently firewalld defualt zone"
// @return        error            error          "Possible errors: ALREADY_ENABLED."
func (c *DbusClientSerivce) PermanentAddPort(port, zone string) (enconterError error) {
	if enconterError = checkPort(port); enconterError == nil {
		if zone == "" {
			zone = c.GetDefaultZone()
		}

		port, protocol := splitPortProtocol(port)

		if path, enconterError := c.generatePath(zone, apis.ZONE_PATH); enconterError == nil {
			obj := c.client.Object(apis.INTERFACE, path)
			printPath(path, apis.CONFIG_ZONE_ADDPORT)
			klog.V(4).Infof("Trying create port Permanent rule in zone %s, %s/%s.", zone, port, protocol)
			call := obj.Call(apis.CONFIG_ZONE_ADDPORT, dbus.FlagNoAutoStart, port, protocol)
			enconterError = call.Err
			if enconterError == nil {
				return nil
			}

		}
	}
	klog.Errorf("Create a permanently port rule failed: %v", enconterError)
	return enconterError
}

/*
 * @title         GetPort
 * @description   temporary get a firewalld port list
 * @auth          author           2021-10-05
 * @param         zone             string         "The empty string is usage default zone, is currently firewalld defualt zone."
 *                                                   e.g. public|dmz..
 * @return        []list           Port           "Returns port list of zone."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE"
 */
func (c *DbusClientSerivce) GetPort(zone string) (list []apis.Port, enconterError error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	call := obj.Call(apis.ZONE_GETPORTS, dbus.FlagNoAutoStart, zone)
	printPath(apis.PATH, apis.ZONE_GETPORTS)
	klog.V(4).Infof("Trying to get port rule in zone %s.", zone)

	enconterError = call.Err
	if enconterError == nil || len(call.Body) >= 0 {
		portList := call.Body[0].([][]string)
		for _, value := range portList {
			list = append(list, apis.Port{
				Port:     value[0],
				Protocol: value[1],
			})
		}
		return list, enconterError
	}
	klog.Errorf("Get a port rule failed: ", enconterError)
	return
}

/*
 * @title         PermanentGetPort
 * @description   get Permanent configurtion a firewalld port list.
 * @auth          author           2021-10-05
 * @param         zone             string         "The empty string is usage default zone, is currently firewalld defualt zone"
 *														e.g. public|dmz..
 * @return        []list           Port           "Returns port list of zone."
 * @return        error            error          "Possible errors:"
 * 														INVALID_ZONE
 */
func (c *DbusClientSerivce) PermanentGetPort(zone string) (list []apis.Port, enconterError error) {

	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, enconterError = c.generatePath(zone, apis.ZONE_PATH); enconterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(apis.ZONE_PATH, apis.CONFIG_ZONE_GETPORTS)
		klog.V(4).Infof("Trying to get permanent port rule in zone %s.", zone)
		call := obj.Call(apis.CONFIG_ZONE_GETPORTS, dbus.FlagNoAutoStart)

		enconterError = call.Err
		if enconterError == nil {
			portList := call.Body[0].([][]interface{})
			for _, value := range portList {
				list = append(list, apis.Port{
					Port:     value[0].(string),
					Protocol: value[1].(string),
				})
			}
		}

	}
	klog.Errorf("Get permanently port Rule failed: %v", enconterError)
	return nil, enconterError
}

/*
 * @title         RemovePort
 * @description   temporary delete a firewalld port
 * @auth      	  author           2021-10-05
 * @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
 * @param         zone    		   string         "e.g. public|dmz.. The empty string is usage default zone, is currently firewalld defualt zone"
 * @return        bool             string         "Returns name of zone from which the port was removed."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE,
 *                                                      INVALID_PORT,
 *                                                      MISSING_PROTOCOL,
 *                                                      INVALID_PROTOCOL,
 *                                                      NOT_ENABLED,
 *                                                      INVALID_COMMAND"
 */
func (c *DbusClientSerivce) RemovePort(port *apis.Port, zone string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_REMOVEPORT)
	klog.V(4).Infof("Trying to remove port rule in zone %s, port rule is: %s/%s", zone, port.Port, port.Protocol)

	call := obj.Call(apis.ZONE_REMOVEPORT, dbus.FlagNoAutoStart, zone, port.Port, port.Protocol)

	if call.Err != nil {
		klog.Errorf("Remove port rule failed: %v", call.Err)
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentRemovePort
 * @description   Permanently delete (port, protocol) from list of ports of zone.
 * @auth      	  author           2021-10-05
 * @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
 * @param         zone    		   string         "The empty string is usage default zone, is currently firewalld defualt zone"
 * 														e.g. public|dmz.."
 * @return        bool             string         "Returns name of zone from which the port was removed."
 * @return        error            error          "Possible errors:
 *                                                      NOT_ENABLED"
 */
func (c *DbusClientSerivce) PermanentRemovePort(port, zone string) (enconterError error) {
	if enconterError = checkPort(port); enconterError == nil {
		if zone == "" {
			zone = c.GetDefaultZone()
		}
		port, protocol := splitPortProtocol(port)

		var path dbus.ObjectPath
		if path, enconterError = c.generatePath(zone, apis.ZONE_PATH); enconterError != nil {
			obj := c.client.Object(apis.INTERFACE, path)

			printPath(path, apis.CONFIG_ZONE_REMOVEPORT)
			klog.V(4).Infof("Try to remove permanent port rule in zone %s, %s/%s.", zone, port, protocol)

			call := obj.Call(apis.CONFIG_ZONE_REMOVEPORT, dbus.FlagNoAutoStart, port, protocol)

			enconterError = call.Err
			if enconterError == nil {
				return nil
			}
		}

	}
	klog.Errorf("remove permanently port rule failed:", enconterError)
	return enconterError
}
