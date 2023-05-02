package firewalld

import (
	"errors"
	"sort"

	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/apis"
)

func (c *DbusClientSerivce) checkZoneName(name string) error {
	if len(name) > 17 {
		klog.Errorf("zone name is limited to 17 chars:", name)
		return errors.New("zone name is limited to 17 chars.")
	}
	return nil
}

func (c *DbusClientSerivce) GetDefaultZone() string {
	return c.defaultZone
}

// @title         SetDefaultZone
// @description   Set default zone for connections and interfaces where no zone has been selected to zone.
// @auther      	  author           2021-09-26
// @param 		  zone			   zone name
// @return        error            error          ""
func (c *DbusClientSerivce) SetDefaultZone(zone string) (err error) {

	//print log
	c.eventLogFormat.Format = ZoneDefaultStartFormat
	c.eventLogFormat.resourceType = "zone"
	c.eventLogFormat.resource = zone
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.INTERFACE_SETDEFAULTZONE)
	call := obj.Call(apis.INTERFACE_SETDEFAULTZONE, dbus.FlagNoAutoStart, zone)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil {
		c.eventLogFormat.Format = ZoneDefaultSuccessFormat
		c.printResourceEventLog()
		return nil
	}
	c.eventLogFormat.Format = ZoneDefaultFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         GetZones
// @description   Return runtime settings of given zone.
// @auther      	  author           2021-09-26
// @return        zones            []string       "Return array of names (s) of predefined zones known to current runtime environment."
// @return        error            error          ""
func (c *DbusClientSerivce) GetZones() ([]string, error) {

	//print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "zone"
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()
	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_GETZONES)
	call := obj.Call(apis.ZONE_GETZONES, dbus.FlagNoAutoStart)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil && len(call.Body) > 0 {
		list := call.Body[0].([]string)
		c.eventLogFormat.Format = ListResourceSuccessFormat
		c.eventLogFormat.resource = list
		c.printResourceEventLog()
		return list, nil
	}

	c.eventLogFormat.Format = ListResourceFailedFormat
	c.printResourceEventLog()
	return nil, call.Err
}

// @title         getZoneId
// @description   Return runtime settings of given zone.
// @auther      	  author           2021-09-26
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
		klog.V(5).Infof("Zone id is: %v", index)
		return index
	} else {
		klog.Warningf("Not Found Zone:", zone)
		return -1
	}
}

// @title         GetZoneSettings
// @description   Return runtime settings of given zone.
// @auther      	  author           2021-09-26
// @param         zone		       string         "zone name."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) GetZoneSettings(zone string) error {

	//print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "zone setting"
	c.eventLogFormat.resource = zone
	c.eventLogFormat.encounterError = nil

	if c.eventLogFormat.encounterError = c.checkZoneName(zone); c.eventLogFormat.encounterError == nil {
		c.printResourceEventLog()
		obj := c.client.Object(apis.INTERFACE, apis.PATH)
		c.printPath(apis.INTERFACE_GETZONESETTINGS)
		call := obj.Call(apis.INTERFACE_GETZONESETTINGS, dbus.FlagNoAutoStart, zone)
		c.eventLogFormat.encounterError = call.Err

		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = ListResourceSuccessFormat
			return nil
		}
	}
	c.eventLogFormat.Format = ListResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         RemoveZone
// @description   Return runtime settings of given zone.
// @auther      	  author           2021-09-26
// @param         zone		       string         "zone name."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) RemoveZone(zone string) error {

	//print log
	c.eventLogFormat.Format = RemoveResourceStartFormat
	c.eventLogFormat.resourceType = "zone"
	c.eventLogFormat.resource = zone
	c.eventLogFormat.encounterError = nil

	if c.eventLogFormat.encounterError = c.checkZoneName(zone); c.eventLogFormat.encounterError == nil {

		path, err := c.generatePath(zone, apis.ZONE_PATH)
		c.eventLogFormat.encounterError = err
		if c.eventLogFormat.encounterError == nil {
			c.printResourceEventLog()
			obj := c.client.Object(apis.INTERFACE, path)

			c.printPath(apis.INTERFACE)
			call := obj.Call(apis.CONFIG_REMOVEZONE, dbus.FlagNoAutoStart)
			c.eventLogFormat.encounterError = call.Err
			if c.eventLogFormat.encounterError == nil {
				c.eventLogFormat.Format = RemoveResourceSuccessFormat
				c.printResourceEventLog()
				return nil
			}
		}
	}
	c.eventLogFormat.Format = RemoveResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         AddZone
// @description   Add zone with given settings into permanent configuration.
// @auther      	  author           2021-09-27
// @param         name		       string         "Is an optional start and end tag and is used to give a more readable name."
// @return        error            error          "Possible errors: NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"
func (c *DbusClientSerivce) AddZone(setting *apis.Settings) error {

	// print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resourceType = "zone"
	c.eventLogFormat.resource = setting.Short
	c.eventLogFormat.encounterError = nil

	if c.eventLogFormat.encounterError = c.checkZoneName(setting.Short); c.eventLogFormat.encounterError == nil {
		c.printResourceEventLog()
		obj := c.client.Object(apis.INTERFACE, apis.CONFIG_PATH)

		c.printPath(apis.CONFIG_ADDZONE)
		call := obj.Call(apis.CONFIG_ADDZONE, dbus.FlagNoAutoStart, setting.Short, setting)
		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = CreateResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}
	c.eventLogFormat.Format = CreateResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         GetZoneOfInterface
// @description   temporary add a firewalld port
// @auther      	  author           2021-09-27
// @param         iface    		   string         "e.g. eth0, iface is device name."
// @return        zoneName         string         "Return name (s) of zone the interface is bound to or empty string.."
func (c *DbusClientSerivce) GetZoneOfInterface(iface string) string {

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "zone"
	c.eventLogFormat.resource = iface
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_GETZONEOFINTERFACE)
	call := obj.Call(apis.ZONE_GETZONEOFINTERFACE, dbus.FlagNoAutoStart, iface)
	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError == nil && len(call.Body) > 0 {
		name, ok := call.Body[0].(string)
		if ok {
			c.eventLogFormat.Format = QueryResourceSuccessFormat
			c.eventLogFormat.resource = name
			c.printResourceEventLog()
			return name
		} else {
			c.eventLogFormat.resource = nil
		}
	}
	c.eventLogFormat.Format = QueryResourceFailedFormat
	c.printResourceEventLog()
	return ""
}

// @title         GetZoneOfInterface
// @description   temporary add a firewalld port
// @auther        author           2023-04-22
// @param         iface    		   string         "e.g. eth0, iface is device name."
// @return        zoneName         string         "Return name (s) of zone the interface is bound to or empty string.."
func (c *DbusClientSerivce) GetDefaultPolicy() string {

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "zone"
	c.eventLogFormat.resource = "get-target"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()
	var path dbus.ObjectPath
	path, c.eventLogFormat.encounterError = c.generatePath(c.defaultZone, apis.ZONE_PATH)
	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printPath(apis.CONFIG_DEFAULT_POLICY)
		call := obj.Call(apis.CONFIG_DEFAULT_POLICY, dbus.FlagNoAutoStart)
		c.eventLogFormat.encounterError = call.Err
		if c.eventLogFormat.encounterError == nil && len(call.Body) > 0 {
			name, ok := call.Body[0].(string)
			if ok {
				c.eventLogFormat.Format = QueryResourceSuccessFormat
				c.eventLogFormat.resource = name
				c.printResourceEventLog()
				return name
			} else {
				c.eventLogFormat.resource = nil
			}
		}
	}
	c.eventLogFormat.Format = QueryResourceFailedFormat
	c.printResourceEventLog()
	return ""
}
