package firewalld

import (
	"github.com/godbus/dbus/v5"

	"github.com/cylonchau/firewalld-gateway/apis"
)

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
func (c *DbusClientSerivce) BindInterface(zone, interfaceName string) (string, error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "NIC"
	c.eventLogFormat.resource = interfaceName
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_ADDINTERFACE)
	call := obj.Call(apis.ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, zone, interfaceName)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError != nil {
		c.eventLogFormat.Format = CreateResourceFailedFormat
		c.printResourceEventLog()
		return "", c.eventLogFormat.encounterError
	}

	c.eventLogFormat.Format = CreateResourceSuccessFormat
	c.printResourceEventLog()
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
func (c *DbusClientSerivce) BindPermanentInterface(zone, interfaceName string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	//print log
	c.eventLogFormat.Format = CreatePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "NIC"
	c.eventLogFormat.resource = interfaceName
	c.eventLogFormat.encounterError = nil

	path, err := c.generatePath(zone, apis.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printPath(apis.CONFIG_ZONE_ADDINTERFACE)
		call := obj.Call(apis.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interfaceName)

		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = CreatePermanentResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}
	c.eventLogFormat.Format = CreatePermanentResourceFailedFormat
	c.printResourceEventLog()
	return err
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
func (c *DbusClientSerivce) QueryInterface(zone, interfaceName string) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "NIC"
	c.eventLogFormat.resource = interfaceName
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_QUERYINTERFACE)
	call := obj.Call(apis.ZONE_QUERYINTERFACE, dbus.FlagNoAutoStart, zone, interfaceName)

	if call.Body[0].(bool) {
		c.eventLogFormat.Format = QueryResourceSuccessFormat
		c.printResourceEventLog()
		return true
	}

	c.eventLogFormat.Format = QueryResourceFailedFormat
	c.printResourceEventLog()
	return false
}

/*
 * @title         PermanentQueryInterface
 * @description   Permanently Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) QueryPermanentInterface(zone, interfaceName string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = QueryPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "NIC"
	c.eventLogFormat.resource = interfaceName
	c.eventLogFormat.encounterError = nil

	path, err := c.generatePath(zone, apis.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		c.printResourceEventLog()
		obj := c.client.Object(apis.INTERFACE, path)

		c.printPath(apis.CONFIG_ZONE_ADDINTERFACE)
		call := obj.Call(apis.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interfaceName)

		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = QueryPermanentResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}
	c.eventLogFormat.Format = QueryPermanentResourceFailedFormat
	c.printResourceEventLog()
	return err
}

/*
 * @title         RemoveInterface
 * @description   Permanently Query whether interface has been bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RemoveInterface(zone, interfaceName string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = RemoveResourceStartFormat
	c.eventLogFormat.resourceType = "NIC"
	c.eventLogFormat.resource = interfaceName
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_REMOVEINTERFACE)
	call := obj.Call(apis.ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, zone, interfaceName)
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
 * @title         PermanentRemoveInterface
 * @description   Permanently remove interface from list of interfaces bound to zone.
 * @auth          author           2021-10-05
 * @param         zone             string         "If zone is empty string, use default zone. e.g. public|dmz..  "
 * @param         interfaceName   string         "interface name. e.g. eth0 | ens33.  "
 * @return        error            error          "Possible errors:
 *                                                       NOT_ENABLED"
 */
func (c *DbusClientSerivce) PermanentRemoveInterface(zone, interfaceName string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	//print log
	c.eventLogFormat.Format = RemovePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "NIC"
	c.eventLogFormat.resource = interfaceName
	c.eventLogFormat.encounterError = nil

	path, err := c.generatePath(zone, apis.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		c.printResourceEventLog()
		obj := c.client.Object(apis.INTERFACE, path)

		c.printPath(apis.CONFIG_ZONE_REMOVEINTERFACE)
		call := obj.Call(apis.CONFIG_ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, interfaceName)
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
