package firewalld

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"sync"

	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
	"github.com/cylonchau/firewalldGateway/config"
	"github.com/cylonchau/firewalldGateway/log"

	"github.com/godbus/dbus/v5"
)

var (
	dbusClient     *DbusClientSerivce
	remotelyBusLck sync.Mutex
	PORT           = "55557"
)

type DbusClientSerivce struct {
	client      *dbus.Conn
	defaultZone string
	ip          string
	port        string
}

func NewDbusClientService(addr string) (*DbusClientSerivce, error) {
	remotelyBusLck.Lock()
	defer remotelyBusLck.Unlock()
	var (
		encounterError error
		conn           *dbus.Conn
	)
	if config.CONFIG.Port == "" {
		klog.V(5).Infof("Start Connect to D-Bus Service: %s:%s", addr, PORT)
	} else {
		PORT = config.CONFIG.Port
		klog.V(5).Infof("Start Connect to D-Bus Service: %s:%s", addr, PORT)
	}

	if dbusClient != nil && dbusClient.client.Connected() {
		return dbusClient, nil
	}
	if conn, encounterError = dbus.Connect("tcp:host="+addr+",port="+PORT, dbus.WithAuth(dbus.AuthAnonymous())); encounterError != nil {
		obj := conn.Object(apis.INTERFACE, apis.PATH)
		call := obj.Call(apis.INTERFACE_GETDEFAULTZONE, dbus.FlagNoAutoStart)
		encounterError = call.Err
		if encounterError == nil {
			return &DbusClientSerivce{
				conn,
				call.Body[0].(string),
				addr,
				PORT,
			}, encounterError
		}
	}
	klog.Errorf("Connect to firewalld client failed:", encounterError)
	return nil, encounterError
}

// @title         Reload
// @description   temporary Add rich language rule into zone.
// @auth      	  author           2021-10-05
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) generatePath(zone, interface_path string) (path dbus.ObjectPath, err error) {
	zoneid := c.getZoneId(zone)
	if zoneid < 0 {
		klog.Errorf("invalid zone:", zone)
		return "", errors.New("invalid zone " + interface_path + zone)
	}
	p := fmt.Sprintf("%s/%d", interface_path, zoneid)
	klog.V(5).Infof("Dbus PATH:", p)
	return dbus.ObjectPath(p), nil
}

func (c *DbusClientSerivce) GetDefaultZone() string {
	return c.defaultZone
}

// @title         SetDefaultZone
// @description   Set default zone for connections and interfaces where no zone has been selected to zone.
// @auth      	  author           2021-09-26
// @param 		  zone			   zone name
// @return        error            error          ""
func (c *DbusClientSerivce) SetDefaultZone(zone string) (err error) {
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	klog.V(5).Infof("Call Remotely Dbus:", apis.PATH, apis.INTERFACE_SETDEFAULTZONE)
	klog.V(4).Infof("Try set default zone to %s", zone)
	var enconnterError error
	currentDefaultZone := c.GetDefaultZone()
	call := obj.Call(apis.INTERFACE_SETDEFAULTZONE, dbus.FlagNoAutoStart, zone)
	enconnterError = call.Err
	if enconnterError == nil {
		klog.V(4).Infof("changed zone %s to %s", currentDefaultZone, zone)
		return nil
	}
	klog.Errorf("set default zone to %s failed: %v", zone, enconnterError)
	return enconnterError
}

// @title         GetZones
// @description   Return runtime settings of given zone.
// @auth      	  author           2021-09-26
// @return        zones            []string       "Return array of names (s) of predefined zones known to current runtime environment."
// @return        error            error          ""
func (c *DbusClientSerivce) GetZones() (zones []string, err error) {
	var enconterError error
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	klog.V(5).Infof("Call Remotely D-Bus: ", apis.PATH, apis.ZONE_GETZONES)
	call := obj.Call(apis.ZONE_GETZONES, dbus.FlagNoAutoStart)
	enconterError = call.Err
	if enconterError == nil || len(call.Body) > 0 {
		klog.V(5).Infof("Get zones: ", call.Body[0])
		return call.Body[0].([]string), nil
	}
	klog.Errorf("Get Zones failed: %v", enconterError)
	return nil, call.Err
}

// @title         getZoneId
// @description   Return runtime settings of given zone.
// @auth      	  author           2021-09-26
// @param         zone		       string         "zone name."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) getZoneId(zone string) int {
	var (
		zoneArray []string
		err       error
	)
	if zoneArray, err = c.GetZones(); err != nil {
		klog.Errorf("Invailed zone id:", zone)
		return -1
	}
	index := sort.SearchStrings(zoneArray, zone)
	if index < len(zoneArray) && zoneArray[index] == zone {
		klog.V(5).Infof("Zone id is: ", index)
		return index
	} else {
		klog.Warningf("Not Found Zone:", zone)
		return -1
	}
}

// @title         GetZoneSettings
// @description   Return runtime settings of given zone.
// @auth      	  author           2021-09-26
// @param         zone		       string         "zone name."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) GetZoneSettings(zone string) (enconterError error) {
	if enconterError = c.checkZoneName(zone); enconterError == nil {
		obj := c.client.Object(apis.INTERFACE, apis.PATH)
		klog.V(5).Infof("Call remote D-Bus:", apis.PATH, apis.INTERFACE_GETZONESETTINGS)
		call := obj.Call(apis.INTERFACE_GETZONESETTINGS, dbus.FlagNoAutoStart, zone)
		enconterError = call.Err
		if enconterError == nil {
			return enconterError
		}
	}
	klog.Errorf("Invailed zone name: %s, error: %v", zone, enconterError)
	return enconterError
}

// @title         RemoveZone
// @description   Return runtime settings of given zone.
// @auth      	  author           2021-09-26
// @param         zone		       string         "zone name."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) RemoveZone(zone string) (err error) {
	if err = c.checkZoneName(zone); err != nil {
		log.Error("invail zone name:", zone)
		return err
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)

	log.Info(fmt.Sprintf("Try to delete zone %s.", zone))
	call := obj.Call(apis.CONFIG_REMOVEZONE, dbus.FlagNoAutoStart)
	if call.Err != nil {
		log.Error("delete zone %s.", zone)
		return call.Err
	}

	return
}

// @title         AddZone
// @description   Add zone with given settings into permanent configuration.
// @auth      	  author           2021-09-27
// @param         name		       string         "Is an optional start and end tag and is used to give a more readable name."
// @return        error            error          "Possible errors: NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"

func (c *DbusClientSerivce) AddZone(setting *apis.Settings) (err error) {
	if err = c.checkZoneName(setting.Short); err != nil {
		return err
	}

	obj := c.client.Object(apis.INTERFACE, apis.CONFIG_PATH)

	log.Debug("Call remote D-Bus:", apis.CONFIG_PATH, apis.CONFIG_ADDZONE)
	log.Info(fmt.Printf("Call ZoneSetting is: %#v", setting))
	call := obj.Call(apis.CONFIG_ADDZONE, dbus.FlagNoAutoStart, setting.Short, setting)

	if call.Err != nil {
		log.Error(fmt.Sprintf("Create ZoneSettiings %s failed:", setting.Short), call.Err.Error())
		return call.Err
	}
	log.Info(fmt.Sprintf("add zoneSetting is: %#v", setting))
	return
}

// @title         GetZoneOfInterface
// @description   temporary add a firewalld port
// @auth      	  author           2021-09-27
// @param         iface    		   string         "e.g. eth0, iface is device name."
// @return        zoneName         string         "Return name (s) of zone the interface is bound to or empty string.."
func (c *DbusClientSerivce) GetZoneOfInterface(iface string) string {
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_GETZONEOFINTERFACE)
	log.Info("Get ZoneOfInterface:", iface)
	call := obj.Call(apis.ZONE_GETZONEOFINTERFACE, dbus.FlagNoAutoStart, iface)
	return call.Body[0].(string)
}

/************************************************** port area ***********************************************************/

// @title         addPort
// @description   temporary add a firewalld port
// @auth      	  author           2021-09-29
// @param         portProtocol     string         "e.g. 80/tcp, 1000-1100/tcp, 80, 1000-1100 default protocol tcp"
// @param         zone    		   string         "e.g. public|dmz.. The empty string is usage default zone, is currently firewalld defualt zone"
// @param         timeout    	   int	          "Timeout, 0 is the permanent effect of the currently service startup state."
// @return        zoneName         string         "Returns name of zone to which the protocol was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_PORT, MISSING_PROTOCOL, INVALID_PROTOCOL, ALREADY_ENABLED, INVALID_COMMAND."

func (c *DbusClientSerivce) AddPort(port *apis.Port, zone string, timeout int) (err error) {

	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_ADDPORT)
	log.Info(fmt.Sprintf("Try To Add Port Rule in Zone %s, %s/%s, life cycle is %d", zone, port.Port, port.Protocol, timeout))
	call := obj.Call(apis.ZONE_ADDPORT, dbus.FlagNoAutoStart, zone, port.Port, port.Protocol, timeout)

	if call.Err != nil || len(call.Body) <= 0 {
		log.Error("Create a Port Rule Failed:", call.Err.Error())
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
func (c *DbusClientSerivce) PermanentAddPort(port, zone string) (err error) {

	if err = checkPort(port); err != nil {
		return err
	}

	if zone == "" {
		zone = c.GetDefaultZone()
	}

	port, protocol := splitPortProtocol(port)

	if path, err := c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	} else {
		obj := c.client.Object(apis.INTERFACE, path)
		log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_ADDPORT)
		log.Info(fmt.Sprintf("Try To Add Port Permanently Rule in Zone %s, %s/%s.", zone, port, protocol))
		call := obj.Call(apis.CONFIG_ZONE_ADDPORT, dbus.FlagNoAutoStart, port, protocol)
		if call.Err != nil {
			log.Error("Create a Port Permanently Rule Failed:", call.Err.Error())
			return call.Err
		}
		return nil
	}
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
func (c *DbusClientSerivce) GetPort(zone string) (list []*apis.Port, err error) {

	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	call := obj.Call(apis.ZONE_GETPORTS, dbus.FlagNoAutoStart, zone)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_GETPORTS)
	log.Info(fmt.Sprintf("Try to get port Rule in Zone %s.", zone))
	if call.Err != nil || len(call.Body) <= 0 {
		log.Error("Get a port Rule failed, not found rule or:", call.Err.Error())
		return nil, call.Err
	}
	portList := call.Body[0].([][]string)
	for _, value := range portList {
		list = append(list, &apis.Port{
			Port:     value[0],
			Protocol: value[1],
		})
	}
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
func (c *DbusClientSerivce) PermanentGetPort(zone string) (list []*apis.Port, err error) {

	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return nil, err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", apis.ZONE_PATH, apis.CONFIG_ZONE_GETPORTS)
	log.Info(fmt.Sprintf("Try to get permanently port Rule in Zone %s.", zone))
	call := obj.Call(apis.CONFIG_ZONE_GETPORTS, dbus.FlagNoAutoStart)

	if call.Err != nil {
		log.Error("get permanently port Rule error:", call.Err.Error())
		return nil, call.Err
	}
	portList := call.Body[0].([][]interface{})

	for _, value := range portList {
		list = append(list, &apis.Port{
			Port:     value[0].(string),
			Protocol: value[1].(string),
		})
	}
	return
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
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_REMOVEPORT)
	log.Info(fmt.Sprintf("Try to delete port Rule in Zone %s, port rule is: %s/%s", zone, port.Port, port.Protocol))
	call := obj.Call(apis.ZONE_REMOVEPORT, dbus.FlagNoAutoStart, zone, port.Port, port.Protocol)

	if call.Err != nil {
		log.Error("remove port rule failed:", call.Err.Error())
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
func (c *DbusClientSerivce) PermanentRemovePort(port, zone string) (b bool, err error) {
	if err = checkPort(port); err != nil {
		return false, err
	}

	if zone == "" {
		zone = c.GetDefaultZone()
	}
	port, protocol := splitPortProtocol(port)

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return false, err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_REMOVEPORT)
	log.Info(fmt.Sprintf("try to remove permanently port rule in Zone %s, %s/%s.", zone, port, protocol))
	call := obj.Call(apis.CONFIG_ZONE_REMOVEPORT, dbus.FlagNoAutoStart, port, protocol)

	if call.Err != nil {
		log.Error("remove permanently port rule failed:", call.Err.Error())
		return false, call.Err
	}
	return true, nil
}

/************************************************** Protocol area ***********************************************************/

// @title         AddProtocol
// @description   temporary get a firewalld port list
// @auth      	  author           2021-09-29
// @param         zone    		   string         "e.g. public|dmz.. If zone is empty string, use default zone. "
// @param         protocol         string         "e.g. tcp|udp... The protocol can be any protocol supported by the system."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the protocol was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_PROTOCOL, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddProtocol(zone, protocol string, timeout int) (list string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_ADDPROTOCOL)
	log.Info(fmt.Sprintf("try to add protocol rule in Zone %s, timeout is %s.", zone, protocol))
	call := obj.Call(apis.ZONE_ADDPROTOCOL, dbus.FlagNoAutoStart, zone, protocol, timeout)

	if call.Err != nil {
		log.Error("add protocol rule failed:", call.Err.Error())
		return "", call.Err
	}
	return call.Body[0].(string), nil
}

/************************************************** service area ***********************************************************/

// @title         NewService
// @description   create new service with given settings into permanent configuration.
// @auth      	  author           2021-10-23
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the service was added."
// @return        error            error          "Possible errors:
//													INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) ListServices() (list []string, err error) {

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.INTERFACE_LISTSERVICES)
	log.Info(fmt.Sprintf("try to list of available services in %s.", c.ip))
	call := obj.Call(apis.INTERFACE_LISTSERVICES, dbus.FlagNoAutoStart)

	if call.Err != nil {
		log.Error("list of available services failed:", call.Err.Error())
		return nil, call.Err
	}
	services := call.Body[0].([]string)
	log.Info(fmt.Sprintf("available services in %s is", c.ip), services)
	return services, nil
}

// @title         NewService
// @description   in runtime configuration.
// @auth      	  author           2021-10-23
// @param         service    	   string         		"service name."
// @param         setting          *ServiceSetting      "service configruate"
// @return        error            error          		"Possible errors:
//															NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"
func (c *DbusClientSerivce) NewService(name string, setting *apis.ServiceSetting) (err error) {

	obj := c.client.Object(apis.INTERFACE, apis.CONFIG_PATH)
	log.Debug("Call remote D-Bus:", apis.CONFIG_PATH, apis.CONFIG_ADDSERVICE)
	log.Info(fmt.Sprintf("try to create a new service in %s.", c.ip))
	log.Debug(fmt.Sprintf("service setting is: %+v", setting))
	call := obj.Call(apis.CONFIG_ADDSERVICE, dbus.FlagNoAutoStart, name, &setting)

	if call.Err != nil {
		log.Error(fmt.Sprintf("create a new service %s failed:.", name), call.Err.Error())
		return call.Err
	}
	log.Info(fmt.Sprintf("create a new service %s success.", name))
	return nil
}

// @title         AddService
// @description   temporary Add service into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the service was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddService(zone, service string, timeout int) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_ADDSERVICE)
	log.Info(fmt.Sprintf("try to add serivce rule in %s Zone %s, timeout is %s.", service, zone, timeout))
	call := obj.Call(apis.ZONE_ADDSERVICE, dbus.FlagNoAutoStart, zone, service, timeout)

	var incurredError error
	incurredError = call.Err
	if incurredError == nil {
		return nil
	}

	log.Error("add service failed:", call.Err.Error())
	return incurredError
}

// @title         PermanentAddService
// @description   Permanent Add service into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentAddService(zone, service string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_ADDSERVICE)
	log.Info(fmt.Sprintf("try to add permanent serivce rule %s in %s Zone %s.", service, zone))
	call := obj.Call(apis.CONFIG_ZONE_ADDSERVICE, dbus.FlagNoAutoStart, service)

	if call.Err != nil {
		log.Error("add permanent service rule failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

// @title         QueryService
// @description   temporary check whether service has been added for zone..
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) QueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_QUERYSERVICE)
	log.Info(fmt.Sprintf("try to query serivce rule %s in %s Zone %s.", service, zone))
	call := obj.Call(apis.ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, zone, service)
	if !call.Body[0].(bool) {
		log.Error("Can not found serivce rule:", service)
		return false
	}
	return true
}

// @title         PermanentQueryService
// @description   Permanent Return whether Add service in rich rules in zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentQueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	var err error
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return false
	}

	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.CONFIG_ZONE_QUERYSERVICE)
	log.Info(fmt.Sprintf("try to query permanent serivce rule %s in %s Zone %s.", service, zone))
	call := obj.Call(apis.CONFIG_ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, service)
	if !call.Body[0].(bool) {
		log.Error("Can not found permanent service rule:", service)
		return false
	}
	return true
}

// @title         RemoveService
// @description   temporary Remove service from zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveService(zone, service string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_REMOVESERVICE)
	log.Info(fmt.Sprintf("try to remove serivce rule %s in %s Zone %s.", service, zone))
	call := obj.Call(apis.ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, zone, service)

	if call.Err != nil {
		log.Error("remove service rule failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

// @title         PermanentAddService
// @description   Permanent Add service into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentRemoveService(zone, service string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_REMOVESERVICE)
	log.Info(fmt.Sprintf("try to remove permanent serivce rule %s in %s Zone %s.", service, zone))
	call := obj.Call(apis.CONFIG_ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, service)

	if call.Err != nil {
		log.Error("remove permanent serivce rule failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

// @title         GetService
// @description   Permanent get service in zone.
// @auth      	  author           2021-10-21
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @return        error            error          "Possible errors:
//														INVALID_ZONE"
func (c *DbusClientSerivce) GetService(zone string) (services []string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_GETSERVICES)
	log.Info(fmt.Sprintf("try to get serivces in zone %s.", zone))
	call := obj.Call(apis.ZONE_GETSERVICES, dbus.FlagNoAutoStart, zone)

	if call.Err != nil {
		log.Error(fmt.Sprintf("get serivces in zone %s failed:", zone), call.Err.Error())
		return nil, call.Err
	}

	services = call.Body[0].([]string)
	log.Info(fmt.Sprintf("serivces in zone %s is:", zone), services)
	return services, nil
}

// @title         PermanentGetServices
// @description   get permanently service in zone.
// @auth      	  author           2021-10-21
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentGetServices(zone, service string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_GETSERVICES)
	log.Info(fmt.Sprintf("try to get permanently serivces in zone %s.", zone))
	call := obj.Call(apis.CONFIG_ZONE_GETSERVICES, dbus.FlagNoAutoStart, service)

	if call.Err != nil {
		log.Error(fmt.Sprintf("get permanently serivces in zone %s failed:", zone), call.Err.Error())
		return call.Err
	}
	return nil
}

/************************************************** Masquerade area ***********************************************************/

/*
 * @title         EnableMasquerade
 * @description   temporary enable masquerade in zone..
 * @auth          author           2021-09-29
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         timeout          int            "Timeout, If timeout is non-zero, masquerading will be active for the amount of seconds."
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE,
 *                                                  ALREADY_ENABLED,
 *                                                  INVALID_COMMAND"
 */
func (c *DbusClientSerivce) EnableMasquerade(zone string, timeout int) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_ADDMASQUERADE)
	log.Info(fmt.Sprintf("try to enable network masquerade in zone %s .", zone))
	call := obj.Call(apis.ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart, zone, timeout)

	if call.Err != nil && len(call.Body) <= 0 {
		log.Error("enable network masquerade failed:", call.Err.Error())
		return call.Err
	}
	return
}

/*
 * @title         PermanentEnableMasquerade
 * @description   permanent enable masquerade in zone..
 * @auth          author           2021-09-29
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE,
 *                                                  ALREADY_ENABLED,
 *                                                  INVALID_COMMAND"
 */
func (c *DbusClientSerivce) PermanentEnableMasquerade(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_ADDMASQUERADE)
	log.Info(fmt.Sprintf("try to permanent enable network masquerade in zone %s .", zone))
	call := obj.Call(apis.CONFIG_ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart)

	if call.Err != nil {
		log.Error("permanent enable network masquerade failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

/*
 * @title         DisableMasquerade
 * @description   temporary enable masquerade in zone..
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         timeout          int            "Timeout, If timeout is non-zero, masquerading will be active for the amount of seconds."
 * @return        zoneName         string         "Returns name of zone in which the masquerade was enabled."
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE,
 *                                                  NOT_ENABLED,
 *                                                  INVALID_COMMAND"
 */
func (c *DbusClientSerivce) DisableMasquerade(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_REMOVEMASQUERADE)
	log.Info(fmt.Sprintf("try to disable network masquerade in zone %s .", zone))
	call := obj.Call(apis.ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart, zone)

	if call.Err != nil && len(call.Body) <= 0 {
		log.Error("disable network masquerade failed:", call.Err.Error())
		return call.Err
	}
	return
}

/*
 * @title         PermanentDisableMasquerade
 * @description   permanent enable masquerade in zone..
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        b            	   bool           "Possible errors:
 * @return        error            error          "Possible errors:
 *                                                  NOT_ENABLED"
 */
func (c *DbusClientSerivce) PermanentDisableMasquerade(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_REMOVEMASQUERADE)
	log.Info(fmt.Sprintf("try to disable network masquerade in zone %s .", zone))

	call := obj.Call(apis.CONFIG_ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart)

	if call.Err != nil && len(call.Body) <= 0 {
		log.Error("disable network masquerade failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentQueryMasquerade
 * @description   query runtime masquerading has been enabled in zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        b            	   bool           "enable: true, disable:false:
 * @return        error            error          "Possible errors:
 *                                                   INVALID_ZONE"
 */
func (c *DbusClientSerivce) PermanentQueryMasquerade(zone string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return false, err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_QUERYMASQUERADE)
	log.Info(fmt.Sprintf("try to query permanent network masquerade in zone %s .", zone))
	call := obj.Call(apis.CONFIG_ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart)

	if call.Err != nil {
		log.Error("query permanent network masquerade is failed:", call.Err)
		return false, call.Err
	}

	if call.Body[0].(bool) == false {
		log.Error("query permanent network masquerade is disabled.")
		return false, nil
	}
	return true, nil
}

/*
 * @title         QueryMasquerade
 * @description   query runtime masquerading has been enabled in zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         timeout          int            "Timeout, If timeout is non-zero, masquerading will be active for the amount of seconds."
 * @return        zoneName         string         "Returns name of zone in which the masquerade was enabled."
 * @return        error            error          "Possible errors:
 *                                                  INVALID_ZONE"
 */
func (c *DbusClientSerivce) QueryMasquerade(zone string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_QUERYMASQUERADE)
	log.Info(fmt.Sprintf("try to query permanent network masquerade in zone %s .", zone))
	call := obj.Call(apis.ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart, zone)
	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		log.Error("query network masquerade is disabled.")
		return false, call.Err
	}
	return true, nil
}

/************************************************** Interface area ***********************************************************/

/*
 * @title         BindInterface
 * @description   temporary Bind interface with zone. From now on all traffic
 * 				   going through the interface will respect the zone's settings.
 * @auth          author           2021-09-29
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        zoneName         string         "Returns name of zone to which the interface was bound."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE,
 *                                                      INVALID_INTERFACE,
 *                                                      ALREADY_ENABLED,
 *                                                      INVALID_COMMAND"
 */
func (c *DbusClientSerivce) BindInterface(zone, interface_name string) (list string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_ADDINTERFACE)
	log.Info(fmt.Sprintf("try to bind interface %s to rule in zone %s .", interface_name, zone))
	call := obj.Call(apis.ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)

	if call.Err != nil {
		log.Error(fmt.Sprintf("bind interface %s failed:", interface_name), call.Err.Error())
		return "", call.Err
	}
	return call.Body[0].(string), nil
}

/*
 * @title         PermanentBindInterface
 * @description   Permanently Bind interface with zone. From now on all traffic
 * 				   going through the interface will respect the zone's settings.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentBindInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_ADDINTERFACE)
	log.Info(fmt.Sprintf("try to permanent bind interface %s to rule in zone %s .", interface_name, zone))
	call := obj.Call(apis.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interface_name)

	if call.Err != nil {
		log.Error(fmt.Sprintf("permanent bind interface %s. failed", interface_name))
		return call.Err
	}
	return nil
}

/*
 * @title         QueryInterface
 * @description   temporary Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         interface        string         "device name， e.g. "
 * @return        b         	   bool           "true:enable, fales:disable."
 * @return        error            error          "Possible errors:
 *                                                      INVALID_ZONE,
 *                                                      INVALID_INTERFACE
 */
func (c *DbusClientSerivce) QueryInterface(zone, interface_name string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_QUERYINTERFACE)
	log.Info(fmt.Sprintf("try to query interface %s bind in zone %s .", interface_name, zone))
	call := obj.Call(apis.ZONE_QUERYINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)

	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		log.Error(fmt.Sprintf("query interface %s bind failed.", interface_name))
		return false, call.Err
	}
	return true, nil
}

/*
 * @title         PermanentQueryInterface
 * @description   Permanently Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) PermanentQueryInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_ADDINTERFACE)
	log.Info(fmt.Sprintf("try to query permanent interface %s bind in zone %s .", interface_name, zone))
	call := obj.Call(apis.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interface_name)

	if call.Err != nil {
		log.Error(fmt.Sprintf("query permanent interface %s bind failed.", interface_name))
		return call.Err
	}
	return nil
}

/*
 * @title         RemoveInterface
 * @description   Permanently Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RemoveInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_REMOVEINTERFACE)
	log.Info(fmt.Sprintf("try to remove interface %s bind in zone %s.", interface_name, zone))
	call := obj.Call(apis.ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)
	if call.Err != nil && len(call.Body) <= 0 {
		log.Error(fmt.Sprintf("remove interface %s bind failed.", interface_name))
		return call.Err
	}
	return nil
}

/*
 * @title         PermanentRemoveInterface
 * @description   Permanently remove interface from list of interfaces bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         interface_name   string         "interface name. e.g. eth0 | ens33.  "
 * @return        error            error          "Possible errors:
 *                                                       NOT_ENABLED"
 */
func (c *DbusClientSerivce) PermanentRemoveInterface(zone, interface_name string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err != nil {
		return err
	}

	obj := c.client.Object(apis.INTERFACE, path)
	log.Debug("Call remote D-Bus:", path, apis.CONFIG_ZONE_REMOVEINTERFACE)
	log.Info(fmt.Sprintf("try to remove permanent interface %s bind in zone %s.", interface_name, zone))
	call := obj.Call(apis.CONFIG_ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, interface_name)
	if call.Err != nil {
		log.Info(fmt.Sprintf("remove permanent interface %s bind failed.", interface_name))
		return call.Err
	}
	return nil
}

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
func (c *DbusClientSerivce) GetForwardPort(zone string) (forwards []*apis.ForwardPort, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	log.Debug("Call remote D-Bus:", apis.PATH, apis.ZONE_GETFORWARDPORT)
	log.Info(fmt.Sprintf("try to get ipv4 forward port rule in zone: %s.", zone))
	call := obj.Call(apis.ZONE_GETFORWARDPORT, dbus.FlagNoAutoStart, zone)
	if call.Err != nil && len(call.Body) <= 0 {
		log.Error(fmt.Sprintf("get ipv4 forward port rule in zone %s failed:", zone), call.Err.Error())
		return nil, call.Err
	}

	for _, value := range call.Body[0].([][]string) {
		forword, err := apis.SliceToStruct(value)
		if err != nil {
			log.Error("convert ipv4 forward port string rule to struct rule failed:", err)
			return nil, err
		}
		forwards = append(forwards, &forword)

	}
	return forwards, nil
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
				//	log.Error("convert ipv4 forward port string rule to struct rule failed:", err)
				//	return nil, err
				//}
				//forwards = append(forwards, forword)
			}
			return forwards, nil
		}
	}
	log.Error(fmt.Sprintf("add forward port in zone %s failed:", zone), encounterError)
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
		log.Error(fmt.Sprintf("remove forward port %s to %s:%s at runtime zone failed:", forword.Port, forword.Protocol, forword.ToPort), call.Err.Error())
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

/************************************************** rich rule area ***********************************************************/

// @title         GetRichRules
// @description   Get list of rich-language rules in zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @return        zoneName         string         "Returns name of zone to which the interface was bound."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) GetRichRules(zone string) (ruleList []*apis.Rule, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_GETRICHRULES)
	klog.V(4).Infof("Try to get all rich rule in zone %s.", zone)
	call := obj.Call(apis.ZONE_GETRICHRULES, dbus.FlagNoAutoStart, zone)

	if call.Err != nil {
		klog.Errorf("Cannot get rich rules in zone %s:", zone, call.Err)
		return nil, call.Err
	}
	for _, value := range call.Body[0].([]string) {
		ruleList = append(ruleList, apis.StringToRule(value))
	}
	klog.V(5).Infof("rich rules: %v", ruleList)
	return
}

// @title         AddRichRule
// @description   temporary Add rich language rule into zone.
// @auth      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) AddRichRule(zone string, rule *apis.Rule, timeout int) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_ADDRICHRULE)
	klog.V(4).Infof("Try to add rich rule in zone %s.", rule.ToString(), zone)
	call := obj.Call(apis.ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString(), timeout)

	if call.Err != nil {
		klog.Errorf("Add rich rule failed:", call.Err)
		return call.Err
	}
	klog.V(5).Infof("Add rich %s rule in zone %s success.", rule.ToString(), zone)
	return nil
}

// @title         PermanentAddRichRule
// @description   Permanently Add rich language rule into zone.
// @auth      	  author           2021-10-05
// @param         zone    	       sting 		  "If zone is empty string, use default zone. e.g. public|dmz..  ""
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) PermanentAddRichRule(zone string, rule *apis.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	path, err := c.generatePath(zone, apis.ZONE_PATH)
	if err != nil {
		return err
	}
	obj := c.client.Object(apis.INTERFACE, path)
	printPath(path, apis.CONFIG_ZONE_ADDRICHRULE)
	klog.V(4).Infof("Try to add permanent rich rule %s in zone %s.", rule.ToString(), zone)

	call := obj.Call(apis.CONFIG_ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, rule.ToString())
	if call.Err != nil {
		log.Error("add permanent rich rule failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

// @title         RemoveRichRule
// @description   temporary Remove rich rule from zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_RULE, NOT_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveRichRule(zone string, rule *apis.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_REOMVERICHRULE)
	klog.V(4).Infof("Try to remove rich rule %s in zone %s.", rule.ToString(), zone)

	call := obj.Call(apis.ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())
	if call.Err != nil {
		klog.Errorf("remove rich rule failed: %v", call.Err)
		return call.Err
	}
	return nil
}

// @title         PermanentAddRichRule
// @description   Permanently Add rich language rule into zone.
// @auth      	  author           2021-10-05
// @param         zone    	       sting 		  "If zone is empty string, use default zone. e.g. public|dmz..  ""
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) PermanentRemoveRichRule(zone string, rule *apis.Rule) (encounterError error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var path dbus.ObjectPath
	if path, encounterError = c.generatePath(zone, apis.ZONE_PATH); encounterError != nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(apis.PATH, apis.CONFIG_ZONE_REOMVERICHRULE)
		klog.V(4).Infof("Try to remove permanent rich rule %s in zone %s.", rule.ToString(), zone)
		call := obj.Call(apis.CONFIG_ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart)
		if encounterError = call.Err; encounterError == nil {
			return nil
		}
	}
	klog.Errorf("remove rich rule %s in zone %s failed:", rule.ToString(), zone, encounterError)
	return encounterError
}

// @title         PermanentQueryRichRule
// @description   Check Permanent Configurtion whether rich rule rule has been added in zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        bool             bool           "Possible errors: INVALID_ZONE, INVALID_RULE"
func (c *DbusClientSerivce) PermanentQueryRichRule(zone string, rule *apis.Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	var (
		path           dbus.ObjectPath
		encounterError error
	)

	if path, encounterError = c.generatePath(zone, apis.ZONE_PATH); encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_QUERYRICHRULE)
		klog.V(4).Infof("Try to query permanent rich rule %s in zone %s.", rule.ToString(), zone)
		call := obj.Call(apis.CONFIG_ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, rule.ToString())
		encounterError = call.Err
		if len(call.Body) == 0 || call.Body[0].(bool) {
			return true
		}
	}
	log.Warnf("Query permanent rich rule %s in zone %s failed: %v", rule.ToString(), zone, encounterError)
	return false
}

// @title         QueryRichRule
// @description   Check whether rich rule is already has.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        bool             bool           "Possible errors: INVALID_ZONE, INVALID_RULE"
func (c *DbusClientSerivce) QueryRichRule(zone string, rule *apis.Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	println(apis.PATH, apis.ZONE_QUERYRICHRULE)
	klog.V(4).Infof("Try to query rich rule %s in zone %s.", rule.ToString, zone)

	call := obj.Call(apis.ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())
	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		log.Warnf("query rich rule %s in zone %s failed: %v", rule.ToString(), zone, call.Err.Error())
		return false
	}
	return true
}

/************************************************** fw service area ***********************************************************/

/*
 * @title         Reload
 * @description   temporary Add rich language rule into zone.
 * @auth          author           2021-10-05
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) Reload() error {
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	println(apis.PATH, apis.INTERFACE_RELOAD)
	klog.V(4).Infof("Try to reload firewalld runtime.")
	call := obj.Call(apis.INTERFACE_RELOAD, dbus.FlagNoAutoStart)

	if call.Err != nil {
		log.Info("reload firewalld failed:", call.Err.Error())
		return call.Err
	}
	return nil
}

/*
 * @title         flush currently zone zoneSettings to default zoneSettings.
 * @description   temporary Add rich language rule into zone.
 * @auth          author           2021-10-05
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RuntimeFlush(zone string) (encounterError error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	defaultZoneSetting := apis.Settings{
		Target:      "default",
		Description: "reset by " + config.CONFIG.AppName,
		Short:       "public",
		Interface:   nil,
		Service: []string{
			"ssh",
			"dhcpv6-client",
		},
		Port: []*apis.Port{
			&apis.Port{
				Port:     config.CONFIG.Dbus_Port,
				Protocol: "tcp",
			},
		},
	}

	var path dbus.ObjectPath
	if path, encounterError = c.generatePath(zone, apis.ZONE_PATH); encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_UPDATE)
		klog.V(4).Infof("Try to flush current active zone (%s).", zone)
		call := obj.Call(apis.CONFIG_UPDATE, dbus.FlagNoAutoStart, defaultZoneSetting)
		encounterError = call.Err
		if encounterError == nil || len(call.Body) <= 0 {
			return nil
		}
	}

	klog.Errorf("Flush current zone (%s) failed: %v", zone, encounterError)
	return encounterError
}

/*
 * @title         destroy
 * @description   off firewalld connection.
 * @auth          author    2021-10-31
 */
func (c *DbusClientSerivce) Destroy() {
	err := c.client.Close()
	if err != nil {
		klog.Errorf("Close D-Bus connection failed, %v", err)
	}
}

func printPath(pathName dbus.ObjectPath, interfaceName string) {
	klog.V(5).Infof("Call remote D-Bus:", pathName, interfaceName)
}
