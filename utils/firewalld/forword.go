package firewalld

import (
	"fmt"
	"net"

	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
)

/************************************************** ForwardPort area ***********************************************************/

/*
 * @title         GetForwardPort
 * @description   temporary get IPv4 forward port in zone.
 * @auth      	  author           2021-10-27
 * @param         zone    		   string         	"If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        forwardPort      set          	"Return array of IPv4 forward ports previously added into zone.
 * @return        error            error          	"Possible errors:
 * 														INVALID_ZONE
 */
func (c *DbusClientSerivce) GetForwardPort(zone string) (forwards []apis.ForwardPort, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_GETFORWARDPORT)
	klog.V(4).Infof("Trying to get ipv4 forward port rule in zone: %s.", zone)

	call := obj.Call(apis.ZONE_GETFORWARDPORT, dbus.FlagNoAutoStart, zone)
	err = call.Err
	if err == nil && len(call.Body) >= 0 {
		for _, value := range call.Body[0].([][]string) {
			forword, err := apis.SliceToStruct(value)
			if err == nil {
				forwards = append(forwards, forword)
				return forwards, nil
			}
		}
	}
	klog.Errorf("Get ipv4 forward port rule in zone %s failed: %v", zone, err)
	return nil, err

}

/*
 * @title         PermanentGetForwardPort
 * @description   permanent get IPv4 forward port in zone.
 * @auth      	  author           2021-10-29
 * @param         zone    		   string         	"If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        forwardPort      set          	"Return array of IPv4 forward ports previously added into zone.
 * @return        error            error          	"Possible errors:
 * 														INVALID_ZONE
 */
func (c *DbusClientSerivce) PermanentGetForwardPort(zone string) ([]apis.ForwardPort, error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var forwards []apis.ForwardPort
	path, encounterError := c.generatePath(zone, apis.ZONE_PATH)
	if encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_GETFORWARDPORT)
		klog.V(4).Infof("Try to get forward port rule in zone: %s.", zone)

		call := obj.Call(apis.CONFIG_GETFORWARDPORT, dbus.FlagNoAutoStart)
		encounterError = call.Err
		if encounterError == nil && len(call.Body) >= 0 {
			for _, value := range call.Body[0].([][]interface{}) {
				fmt.Printf("%+v\n", value)
				//forword, err := apis.SliceToStruct(value)
				//if err != nil {
				//	klog.Errorf("convert ipv4 forward port string rule to struct rule failed: %v", err)
				//	return nil, err
				//}
				//forwards = append(forwards, forword)
			}
			return forwards, nil
		}
	}
	klog.Errorf("add forward port in zone %s failed: %v", zone, encounterError)
	return forwards, encounterError
}

/*
 * @title         AddForwardPort
 * @description   temporary Add the IPv4 forward port into zone.
 * @auth      	  author           2021-09-29
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
func (c *DbusClientSerivce) AddForwardPort(zone string, timeout int, forward *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	printPath(apis.PATH, apis.ZONE_ADDFORWARDPORT)
	klog.V(4).Infof("Try to add forward port %s to %s:%s.", forward.Port, forward.ToAddr, forward.ToPort)

	call := obj.Call(apis.ZONE_ADDFORWARDPORT, dbus.FlagNoAutoStart, zone, forward.Port, forward.Protocol, forward.ToPort, forward.ToAddr, timeout)
	if call.Err != nil && len(call.Body) <= 0 {
		klog.Errorf("Add forward port in zone %s failed: %v", zone, call.Err)
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentAddForwardPort
 * @description   temporary Add the IPv4 forward port into zone.
 * @auth      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentAddForwardPort(zone string, forwardPort *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	path, encounterError := c.generatePath(zone, apis.ZONE_PATH)
	if encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)

		printPath(path, apis.CONFIG_ZONE_ADDFORWARDPORT)
		klog.V(4).Infof("try to add permanent forward port %s to %s:%s.", forwardPort.Port, forwardPort.ToAddr, forwardPort.ToPort)

		call := obj.Call(apis.CONFIG_ZONE_ADDFORWARDPORT, dbus.FlagNoAutoStart, forwardPort.Port, forwardPort.Protocol, forwardPort.ToPort, forwardPort.ToAddr)
		encounterError = call.Err
		if encounterError == nil && len(call.Body) >= 0 {
			return nil
		}
	}
	klog.Errorf("Add permanent forward port in zone %s failed: %v", zone, encounterError)
	return encounterError
}

/*
 * @title         RemoveForwardPort
 * @description   temporary (runtime) Remove IPv4 forward port ((port, protocol, toport, toaddr)) from zone.
 * @auth      	  author           2021-09-29
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
func (c *DbusClientSerivce) RemoveForwardPort(zone string, forword *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_REMOVEFORWARDPORT)
	klog.V(4).Infof("Try to remove forward port %s to %s:%s.", forword.Port, forword.ToAddr, forword.ToPort)

	call := obj.Call(apis.ZONE_REMOVEFORWARDPORT, dbus.FlagNoAutoStart, zone, forword.Port, forword.Protocol, forword.ToPort, forword.ToAddr)
	if call.Err != nil && len(call.Body) <= 0 {
		klog.Errorf("remove forward port %s to %s:%s at runtime zone failed: %v", forword.Port, forword.Protocol, forword.ToPort, call.Err.Error())
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentRemoveForwardPort
 * @description   Permanently remove (port, protocol, toport, toaddr) from list of forward ports of zone.
 * @auth      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentRemoveForwardPort(zone string, forword *apis.ForwardPort) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	path, enconterError := c.generatePath(zone, apis.ZONE_PATH)
	if enconterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_REMOVEFORWARDPORT)
		klog.V(4).Infof("Try to remove permanent forward port %s to %s:%s.", forword.Port, forword.ToAddr, forword.ToPort)
		call := obj.Call(apis.CONFIG_ZONE_REMOVEFORWARDPORT, dbus.FlagNoAutoStart, forword.Port, forword.Protocol, forword.ToPort, forword.ToAddr)
		enconterError = call.Err
		if enconterError == nil && len(call.Body) >= 0 {
			return nil
		}
	}
	klog.Errorf("Try to remove permanent forward port  %s to %s:%s failed: %v", forword.Port, forword.ToAddr, forword.ToPort, enconterError)
	return enconterError
}

/*
 * @title         QueryForwardPort
 * @description   temporary (runtime) query whether the IPv4 forward port (port, protocol, toport, toaddr) has been added into zone.
 * @auth      	  author           2021-10-07
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
func (c *DbusClientSerivce) QueryForwardPort(zone string, portProtocol, toHostPort string) (b bool) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var enconterError error
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, enconterError := net.SplitHostPort(toHostPort)
	if enconterError == nil {
		obj := c.client.Object(apis.INTERFACE, apis.PATH)
		printPath(apis.PATH, apis.ZONE_QUERYFORWARDPORT)
		klog.V(4).Infof("Try to query forward port %s to %s.", portProtocol, toHostPort)
		call := obj.Call(apis.ZONE_QUERYFORWARDPORT, dbus.FlagNoAutoStart, zone, port, protocol, toPort, toAddr)
		enconterError = call.Err
		if enconterError == nil || call.Body[0].(bool) {
			return true
		}
	}
	klog.Errorf("Query forward port %s to %s failed: %v", portProtocol, toHostPort, enconterError)
	return false
}

/*
 * @title         PermanentQueryForwardPort
 * @description   Permanently remove (port, protocol, toport, toaddr) from list of forward ports of zone.
 * @auth      	  author           2021-10-07
 * @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         portProtocol     string         "The port can either be a single port number portid or a port
 *													range portid-portid. The protocol can either be tcp or udp e.g. 10-20/tcp|20|20/tcp"
 * @param         toHostPort       string		  "The destination address is a simple IP address. e.g. 10.0.0.1:22"
 * @return        error            error          "Possible errors:
 * 													ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentQueryForwardPort(zone string, portProtocol, toHostPort string) (b bool) {
	var enconterError error
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(portProtocol)
	toAddr, toPort, enconterError := net.SplitHostPort(toHostPort)
	if enconterError == nil {
		var path dbus.ObjectPath
		if path, enconterError = c.generatePath(zone, apis.ZONE_PATH); enconterError == nil {
			obj := c.client.Object(apis.INTERFACE, path)
			printPath(path, apis.CONFIG_ZONE_QUERYFORWARDPORT)
			klog.V(4).Infof("Try to query permanent forward port %s to %s.", portProtocol, toHostPort)
			call := obj.Call(apis.CONFIG_ZONE_QUERYFORWARDPORT, dbus.FlagNoAutoStart, port, protocol, toPort, toAddr)
			enconterError = call.Err
			if enconterError == nil || (len(call.Body) >= 0 || call.Body[0].(bool)) {
				return true
			}
		}
	}
	klog.Errorf("Query permanent forward port %s to %s failed: %v", portProtocol, toHostPort, enconterError)
	return false
}
