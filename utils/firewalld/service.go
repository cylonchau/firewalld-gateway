package firewalld

import (
	"errors"

	"github.com/godbus/dbus/v5"

	"github.com/cylonchau/firewalld-gateway/apis"
)

/************************************************** service area ***********************************************************/

// @title         NewService
// @description   create new service with given settings into permanent configuration.
// @middlewares      	  author           2021-10-23
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the service was added."
// @return        error            error          "Possible errors:
//													INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) GetServices() (list []string, err error) {
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	c.printPath(apis.INTERFACE_LISTSERVICES)

	//print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	call := obj.Call(apis.INTERFACE_LISTSERVICES, dbus.FlagNoAutoStart)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil {
		services, ok := call.Body[0].([]string)
		if ok {
			c.eventLogFormat.Format = ListResourceSuccessFormat
			c.eventLogFormat.resource = services
			c.printResourceEventLog()
			return services, nil
		} else {
			c.eventLogFormat.encounterError = errors.New("reflect resource failed")
		}
	}
	c.eventLogFormat.Format = ListResourceFailedFormat
	c.printResourceEventLog()
	return nil, call.Err
}

// @title         NewService
// @description   in runtime configuration.
// @middlewares      	  author           2021-10-23
// @param         service    	   string         		"service name."
// @param         setting          *ServiceSetting      "service configruate"
// @return        error            error          		"Possible errors:
//															NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"
func (c *DbusClientSerivce) AddNewService(name string, setting *apis.ServiceSetting) error {

	obj := c.client.Object(apis.INTERFACE, apis.CONFIG_PATH)

	//print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = name
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()

	c.printPath(apis.CONFIG_ADDSERVICE)
	call := obj.Call(apis.CONFIG_ADDSERVICE, dbus.FlagNoAutoStart, name, &setting)

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

// @title         AddService
// @description   temporary Add service into zone.
// @middlewares      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        zoneName         string         "Returns name of zone to which the service was added."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddService(zone, service string, timeout uint32) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	// print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = service
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	c.printPath(apis.ZONE_ADDSERVICE)
	call := obj.Call(apis.ZONE_ADDSERVICE, dbus.FlagNoAutoStart, zone, service, timeout)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil {
		c.eventLogFormat.Format = CreateResourceSuccessFormat
		c.eventLogFormat.resource = service
		c.printResourceEventLog()
		return nil
	}

	c.eventLogFormat.Format = CreateResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         PermanentAddService
// @description   Permanent Add service into zone.
// @middlewares      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentAddService(zone, service string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = CreatePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, apis.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printResourceEventLog()

		c.printPath(apis.CONFIG_ZONE_ADDSERVICE)
		call := obj.Call(apis.CONFIG_ZONE_ADDSERVICE, dbus.FlagNoAutoStart, service)

		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = CreatePermanentResourceSuccessFormat
			return nil
		}
	}
	c.eventLogFormat.Format = CreatePermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         QueryService
// @description   temporary check whether service has been added for zone..
// @middlewares      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) QueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = service

	c.printResourceEventLog()

	c.printPath(apis.ZONE_QUERYSERVICE)
	call := obj.Call(apis.ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, zone, service)
	if !call.Body[0].(bool) {
		c.eventLogFormat.Format = QueryNotFount
		c.printResourceEventLog()
		return false
	}
	return true
}

// @title         PermanentQueryService
// @description   Permanent Return whether Add service in rich rules in zone.
// @middlewares      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentQueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = QueryPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil

	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, apis.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)

		c.printResourceEventLog()

		c.printPath(apis.CONFIG_ZONE_QUERYSERVICE)
		call := obj.Call(apis.CONFIG_ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, service)

		if call.Body[0].(bool) {
			return true
		} else {
			c.eventLogFormat.encounterError = call.Err
		}
	}
	c.eventLogFormat.Format = QueryPermanentResourceFailedFormat
	c.printResourceEventLog()
	return false
}

// @title         RemoveService
// @description   temporary Remove service from zone.
// @middlewares      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveService(zone, service string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = RemoveResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = service
	c.eventLogFormat.encounterError = nil

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_REMOVESERVICE)
	call := obj.Call(apis.ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, zone, service)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError != nil {
		c.eventLogFormat.Format = RemoveResourceFailedFormat
		c.printResourceEventLog()
		return c.eventLogFormat.encounterError
	}
	c.eventLogFormat.Format = RemoveResourceSuccessFormat
	c.printResourceEventLog()
	return nil
}

// @title         PermanentAddService
// @description   Permanent Add service into zone.
// @middlewares      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         service          string         "service name e.g. http|ssh|ftp.."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemovePermanentService(zone, service string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = RemovePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = service
	c.eventLogFormat.encounterError = nil

	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, apis.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printResourceEventLog()

		c.printPath(apis.CONFIG_ZONE_REMOVESERVICE)
		call := obj.Call(apis.CONFIG_ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, service)
		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.resourceType = RemoveResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}
	c.eventLogFormat.resourceType = RemovePermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         PermanentGetServices
// @description   get permanently service in zone.
// @middlewares      	  author           2021-10-21
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) GetPermanentServices(zone, service string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = RemovePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = service
	c.eventLogFormat.encounterError = nil

	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, apis.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)

		c.printResourceEventLog()

		c.printPath(apis.CONFIG_ZONE_GETSERVICES)
		call := obj.Call(apis.CONFIG_ZONE_GETSERVICES, dbus.FlagNoAutoStart, service)

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
