package firewalld

import (
	"k8s.io/klog/v2"

	"github.com/godbus/dbus/v5"

	"github.com/cylonchau/firewalldGateway/apis"
)

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
	printPath(apis.PATH, apis.INTERFACE_LISTSERVICES)

	klog.V(4).Infof("Trying get list of available services in %s.", c.ip)

	call := obj.Call(apis.INTERFACE_LISTSERVICES, dbus.FlagNoAutoStart)

	if call.Err == nil {
		services := call.Body[0].([]string)
		klog.V(4).Infof("Available services in %s: %v", c.ip, services)
		return services, nil
	}

	klog.Errorf("list of available services failed: %v", call.Err.Error())
	return nil, call.Err
}

// @title         NewService
// @description   in runtime configuration.
// @auth      	  author           2021-10-23
// @param         service    	   string         		"service name."
// @param         setting          *ServiceSetting      "service configruate"
// @return        error            error          		"Possible errors:
//															NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"
func (c *DbusClientSerivce) NewService(name string, setting *apis.ServiceSetting) error {

	obj := c.client.Object(apis.INTERFACE, apis.CONFIG_PATH)
	printPath(apis.CONFIG_PATH, apis.CONFIG_ADDSERVICE)

	klog.V(4).Infof("Trying create a new service in %s.", c.ip)
	klog.V(5).Infof("Service setting is: %+v", setting)

	call := obj.Call(apis.CONFIG_ADDSERVICE, dbus.FlagNoAutoStart, name, &setting)

	if call.Err != nil {
		klog.Errorf("Create a new service %s failed: %v", name, call.Err.Error())
		return call.Err
	}
	klog.V(4).Infof("Create a new service %s success.", name)
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
	printPath(apis.PATH, apis.ZONE_ADDSERVICE)
	klog.V(4).Infof("Trying to create serivce rule in %s zone %s, timeout is %s.", service, zone, timeout)
	call := obj.Call(apis.ZONE_ADDSERVICE, dbus.FlagNoAutoStart, zone, service, timeout)

	var incurredError error
	incurredError = call.Err
	if incurredError == nil {
		return nil
	}

	klog.Errorf("Create service failed: %v", incurredError)
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
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_ADDSERVICE)
		klog.V(4).Infof("Trying to create permanent serivce rule %s in %s zone %s.", service, zone)
		call := obj.Call(apis.CONFIG_ZONE_ADDSERVICE, dbus.FlagNoAutoStart, service)
		err = call.Err
		if call.Err == nil {
			return nil
		}
	}
	klog.Errorf("Create permanent service rule failed: %v", err)
	return err
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

	printPath(apis.PATH, apis.ZONE_QUERYSERVICE)
	klog.V(4).Infof("Trying to query serivce rule %s in %s zone %s.", service, zone)

	call := obj.Call(apis.ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, zone, service)
	if !call.Body[0].(bool) {
		klog.Warningf("Cannot found serivce rule:", service)
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

	var (
		path           dbus.ObjectPath
		encounterError error
	)
	if path, encounterError = c.generatePath(zone, apis.ZONE_PATH); encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(apis.PATH, apis.CONFIG_ZONE_QUERYSERVICE)
		klog.V(4).Infof("try to query permanent serivce rule %s in %s Zone %s.", service, zone)
		call := obj.Call(apis.CONFIG_ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, service)
		if call.Body[0].(bool) {
			return true
		}
	}
	klog.Errorf("Cannot found permanent service rule:", service)
	return false
}

// @title         RemoveService
// @description   temporary Remove service from zone.
// @auth      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveService(zone, service string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_REMOVESERVICE)
	klog.V(4).Infof("Trying to remove serivce rule %s in %s zone %s.", service, zone)

	call := obj.Call(apis.ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, zone, service)
	if call.Err != nil {
		klog.Errorf("Remove service rule failed: %v", call.Err.Error())
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
func (c *DbusClientSerivce) PermanentRemoveService(zone, service string) (enconterError error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	if path, enconterError := c.generatePath(zone, apis.ZONE_PATH); enconterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_REMOVESERVICE)
		klog.V(4).Infof("Trying to remove permanent serivce rule %s in %s zone %s.", service, zone)
		call := obj.Call(apis.CONFIG_ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, service)
		enconterError = call.Err
		if enconterError == nil {
			return nil
		}
	}
	klog.Errorf("Remove permanent serivce rule failed: %v", enconterError)
	return enconterError
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
	printPath(apis.PATH, apis.ZONE_GETSERVICES)
	klog.V(4).Infof("Trying to get serivce in zone %s.", zone)
	call := obj.Call(apis.ZONE_GETSERVICES, dbus.FlagNoAutoStart, zone)

	if call.Err != nil {
		klog.Errorf("get serivces in zone %s failed: %v", zone, call.Err)
		return nil, call.Err
	}

	services = call.Body[0].([]string)
	klog.V(4).Infof("Get serivce in zone %s is %v", zone, services)
	return services, nil
}

// @title         PermanentGetServices
// @description   get permanently service in zone.
// @auth      	  author           2021-10-21
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentGetServices(zone, service string) (encounterError error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	var path dbus.ObjectPath
	if path, encounterError = c.generatePath(zone, apis.ZONE_PATH); encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_GETSERVICES)
		klog.V(4).Infof("Trying to get permanently serivces in zone %s.", zone)
		call := obj.Call(apis.CONFIG_ZONE_GETSERVICES, dbus.FlagNoAutoStart, service)
		encounterError = call.Err
		if encounterError == nil {
			return nil
		}

	}
	klog.Errorf("Get permanently serivces in zone %s failed: %v", zone, encounterError)
	return encounterError
}
