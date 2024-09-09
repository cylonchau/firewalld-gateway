//go:build !swagger
// +build !swagger

package firewalld

import (
	"errors"

	"github.com/godbus/dbus/v5"

	api2 "github.com/cylonchau/firewalld-gateway/api"
)

/************************************************** service area ***********************************************************/

// :title         ListServices
// :description   Return array of avliable of services in permanent configuration.
// :Create        author   2021-10-23
// :Update        author   2024-09-06
// :param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :param         service         string         "service name e.g. http|ssh|ftp.."
// :param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// :return        zoneName        string         "Returns name of zone to which the service was added."
// :return        error           error          "Possible errors:
//
//	INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) ListServices() (list []string, err error) {
	obj := c.client.Object(api2.INTERFACE, api2.PATH)
	c.printPath(api2.INTERFACE_LISTSERVICES)

	//print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	call := obj.Call(api2.INTERFACE_LISTSERVICES, dbus.FlagNoAutoStart)

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

// :title         NewService
// :description   in runtime configuration.
// :Create        author   2021-10-23
// :Update        author   2024-09-06
// :param         service    	 string         		"service name."
// :param         setting        *ServiceSetting      "service configruate"
// :return        error          error          		"Possible errors:
//
//	NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"
func (c *DbusClientSerivce) AddNewService(name string, setting *api2.ServiceSetting) error {

	obj := c.client.Object(api2.INTERFACE, api2.CONFIG_PATH)

	//print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = name
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()

	c.printPath(api2.CONFIG_ADDSERVICE)
	call := obj.Call(api2.CONFIG_ADDSERVICE, dbus.FlagNoAutoStart, name, &setting)

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

// :title         AddServiceRuntime
// :description   Add service into zone (runtime).
// :Create        author   2021-09-29
// :Update        author   2024-09-07
// :param         zone          string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :param         service       string         "service name e.g. http|ssh|ftp.."
// :param         timeout    	int	           "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// :return        zoneName      string         "Returns name of zone to which the service was added."
// :return        error         error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddServiceRuntime(zone, service string, timeout uint32) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(api2.INTERFACE, api2.PATH)

	// print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = service
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	c.printPath(api2.ZONE_ADDSERVICE)
	call := obj.Call(api2.ZONE_ADDSERVICE, dbus.FlagNoAutoStart, zone, service, timeout)

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

// :title         AddPermanentService
// :description   Permanently add service to list of services used in zone.
// :Create        author      2021-09-29
// :Update        author      2024-09-06
// :param         zone        string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :param         service     string         "service name e.g. http|ssh|ftp.."
// :return        error       error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) AddPermanentService(zone, service string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = CreatePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, api2.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(api2.INTERFACE, path)
		c.printResourceEventLog()

		c.printPath(api2.CONFIG_ZONE_ADDSERVICE)
		call := obj.Call(api2.CONFIG_ZONE_ADDSERVICE, dbus.FlagNoAutoStart, service)

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

// :title         QueryService
// :description   temporary check whether service has been added for zone..
// :Create        author   2021-10-05
// :Update        author   2024-09-06
// :param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :param         service          string         "service name e.g. http|ssh|ftp.."
// :return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) QueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(api2.INTERFACE, api2.PATH)

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	c.eventLogFormat.resource = service

	c.printResourceEventLog()

	c.printPath(api2.ZONE_QUERYSERVICE)
	call := obj.Call(api2.ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, zone, service)
	if !call.Body[0].(bool) {
		c.eventLogFormat.Format = QueryNotFount
		c.printResourceEventLog()
		return false
	}
	return true
}

// :title         PermanentQueryService
// :description   Permanent Return whether Add service in rich rules in zone.
// :Create        author   2021-10-05
// :Update        author   2024-09-06
// :param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :param         service          string         "service name e.g. http|ssh|ftp.."
// :return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) PermanentQueryService(zone, service string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = QueryPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil

	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, api2.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(api2.INTERFACE, path)

		c.printResourceEventLog()

		c.printPath(api2.CONFIG_ZONE_QUERYSERVICE)
		call := obj.Call(api2.CONFIG_ZONE_QUERYSERVICE, dbus.FlagNoAutoStart, service)

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

// :title         RemoveService
// :description   Remove service from zone (runtime).
// :Create        author   2021-10-05
// :Update        author   2024-09-07
// :param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :param         service          string         "service name e.g. http|ssh|ftp.."
// :return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveRuntimeService(zone, service string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = RemoveResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = service
	c.eventLogFormat.encounterError = nil

	obj := c.client.Object(api2.INTERFACE, api2.PATH)

	c.printPath(api2.ZONE_REMOVESERVICE)
	call := obj.Call(api2.ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, zone, service)

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

// :title         Permanently remove service from list of services used in zone.
// :description   PPermanently remove service from list of services used in zone.
// :Create        author   2021-09-29
// :Update        author   2024-09-06
// :param         zone     string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :param         service  string         "service name e.g. http|ssh|ftp.."
// :return        error    error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemovePermanentService(service string) error {
	zone := c.GetDefaultZone()

	// print log
	c.eventLogFormat.Format = RemovePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.resource = service
	c.eventLogFormat.encounterError = nil

	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, api2.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(api2.INTERFACE, path)
		c.printResourceEventLog()

		c.printPath(api2.CONFIG_ZONE_REMOVESERVICE)
		call := obj.Call(api2.CONFIG_ZONE_REMOVESERVICE, dbus.FlagNoAutoStart, service)
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

// :title         PermanentGetServices
// :description   Get list of service names used in zone (permanent).
// :Create        author   2021-10-21
// :Update        author   2024-09-06
// :param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :return        error            error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) GetPermanentServices() (list []string, err error) {
	zone := c.GetDefaultZone()

	// print log
	c.eventLogFormat.Format = ListPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	var path dbus.ObjectPath
	if path, c.eventLogFormat.encounterError = c.generatePath(zone, api2.ZONE_PATH); c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(api2.INTERFACE, path)
		c.printResourceEventLog()

		c.printPath(api2.CONFIG_ZONE_GETSERVICES)
		call := obj.Call(api2.CONFIG_ZONE_GETSERVICES, dbus.FlagNoAutoStart)
		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			services, ok := call.Body[0].([]string)
			if ok {
				c.eventLogFormat.Format = ListPermanentResourceSuccessFormat
				c.eventLogFormat.resource = services
				c.printResourceEventLog()
				return services, nil
			} else {
				c.eventLogFormat.encounterError = errors.New("reflect resource failed")
			}
		}
	}

	c.eventLogFormat.Format = ListPermanentResourceFailedFormat
	c.printResourceEventLog()
	return nil, c.eventLogFormat.encounterError
}

// :title         GetRuntimeServices
// :description   Get list of service names used in zone (runtime).
// :Create        author      2024-09-07
// :Update        author      2024-09-07
// :param         zone        string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// :return        error       error          "Possible errors: INVALID_ZONE, INVALID_SERVICE, ALREADY_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) GetRuntimeServices(zone string) (list []string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "service"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(api2.INTERFACE, api2.PATH)
	c.printResourceEventLog()

	c.printPath(api2.ZONE_GETSERVICES)
	call := obj.Call(api2.ZONE_GETSERVICES, dbus.FlagNoAutoStart, zone)
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
	return nil, c.eventLogFormat.encounterError
}
