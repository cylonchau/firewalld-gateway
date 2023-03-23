package firewalld

import (
	"github.com/godbus/dbus/v5"

	"github.com/cylonchau/firewalld-gateway/apis"
)

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
func (c *DbusClientSerivce) EnableMasquerade(zone string, timeout uint32) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	//print log
	c.eventLogFormat.Format = SwitchResourceStartFormat
	c.eventLogFormat.resourceType = "enable"
	c.eventLogFormat.encounterError = nil

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	c.printResourceEventLog()

	c.printPath(apis.ZONE_ADDMASQUERADE)
	call := obj.Call(apis.ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart, zone, timeout)
	c.eventLogFormat.encounterError = call.Err

	if c.eventLogFormat.encounterError == nil {
		c.eventLogFormat.Format = SwitchResourceSuccessFormat
		c.printResourceEventLog()
		return nil
	}

	c.eventLogFormat.Format = SwitchResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
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
func (c *DbusClientSerivce) EnablePermanentMasquerade(zone string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = SwitchPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "enable"
	c.eventLogFormat.encounterError = nil
	path, err := c.generatePath(zone, apis.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printPath(apis.CONFIG_ZONE_ADDMASQUERADE)

		c.printResourceEventLog()

		c.printPath(apis.CONFIG_ZONE_ADDMASQUERADE)
		call := obj.Call(apis.CONFIG_ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart)
		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = SwitchPermanentResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}
	c.eventLogFormat.Format = SwitchPermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
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

	// print log
	c.eventLogFormat.Format = SwitchResourceStartFormat
	c.eventLogFormat.resourceType = "disable"
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_REMOVEMASQUERADE)
	call := obj.Call(apis.ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart, zone)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError != nil {
		c.eventLogFormat.Format = SwitchResourceFailedFormat
		c.printResourceEventLog()
		return c.eventLogFormat.encounterError
	}

	c.eventLogFormat.Format = SwitchResourceSuccessFormat
	c.printResourceEventLog()
	return nil
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
func (c *DbusClientSerivce) DisablePermanentMasquerade(zone string) (err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = SwitchPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "disable"
	c.eventLogFormat.encounterError = nil

	path, err := c.generatePath(zone, apis.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		c.printResourceEventLog()
		obj := c.client.Object(apis.INTERFACE, path)

		c.printPath(apis.CONFIG_ZONE_REMOVEMASQUERADE)
		call := obj.Call(apis.CONFIG_ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart)
		c.eventLogFormat.encounterError = call.Err

		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = SwitchPermanentResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}

	c.eventLogFormat.Format = SwitchPermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
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
func (c *DbusClientSerivce) QueryPermanentMasquerade(zone string) (bool, error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = QueryPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "masquerade"
	c.eventLogFormat.encounterError = nil

	path, err := c.generatePath(zone, apis.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		c.printResourceEventLog()
		obj := c.client.Object(apis.INTERFACE, path)

		c.printPath(apis.CONFIG_ZONE_QUERYMASQUERADE)
		call := obj.Call(apis.CONFIG_ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart)

		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			if call.Body[0].(bool) {
				c.eventLogFormat.Format = QueryPermanentResourceSuccessFormat
				c.printResourceEventLog()
				return true, nil
			}
		}
	}
	c.eventLogFormat.Format = QueryPermanentResourceFailedFormat
	c.printResourceEventLog()
	return false, err
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

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "masquerade"
	c.eventLogFormat.resource = zone
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_QUERYMASQUERADE)
	call := obj.Call(apis.ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart, zone)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil {
		if call.Body[0].(bool) {
			c.eventLogFormat.Format = QueryResourceSuccessFormat
			c.eventLogFormat.resource = "enable"
			c.printResourceEventLog()
			return true, nil
		} else {
			c.eventLogFormat.resource = "disable"
		}
	}
	c.eventLogFormat.Format = QueryResourceFailedFormat
	c.printResourceEventLog()
	return false, c.eventLogFormat.encounterError
}
