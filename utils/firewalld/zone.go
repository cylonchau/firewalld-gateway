package firewalld

import (
	"errors"
	"sort"

	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
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
	klog.V(5).Infof("Call Remotely D-Bus: %s=>%s ", apis.PATH, apis.ZONE_GETZONES)
	call := obj.Call(apis.ZONE_GETZONES, dbus.FlagNoAutoStart)
	enconterError = call.Err
	if enconterError == nil || len(call.Body) > 0 {
		klog.V(5).Infof("Get zones: %v", call.Body[0])
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
		klog.V(5).Infof("Zone id is: %v", index)
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
		printPath(apis.PATH, apis.INTERFACE_GETZONESETTINGS)
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
func (c *DbusClientSerivce) RemoveZone(zone string) (enconterError error) {
	if enconterError = c.checkZoneName(zone); enconterError == nil {
		var path dbus.ObjectPath
		if path, enconterError = c.generatePath(zone, apis.ZONE_PATH); enconterError == nil {
			obj := c.client.Object(apis.INTERFACE, path)
			klog.V(4).Infof("Try to delete zone %s.", zone)
			call := obj.Call(apis.CONFIG_REMOVEZONE, dbus.FlagNoAutoStart)
			enconterError = call.Err
			if enconterError == nil {
				return nil
			}
		}
	}
	klog.Errorf("Delete zone %s failed: %v", zone, enconterError)
	return enconterError
}

// @title         AddZone
// @description   Add zone with given settings into permanent configuration.
// @auth      	  author           2021-09-27
// @param         name		       string         "Is an optional start and end tag and is used to give a more readable name."
// @return        error            error          "Possible errors: NAME_CONFLICT, INVALID_NAME, INVALID_TYPE"
func (c *DbusClientSerivce) AddZone(setting *apis.Settings) (enconterError error) {
	if enconterError = c.checkZoneName(setting.Short); enconterError == nil {
		obj := c.client.Object(apis.INTERFACE, apis.CONFIG_PATH)
		printPath(apis.CONFIG_PATH, apis.CONFIG_ADDZONE)
		klog.V(4).Infof("Call ZoneSetting is: %v", setting)

		call := obj.Call(apis.CONFIG_ADDZONE, dbus.FlagNoAutoStart, setting.Short, setting)
		enconterError = call.Err
		if enconterError == nil {
			klog.V(4).Infof("Add zoneSetting is: %v", setting)
			return nil
		}
	}
	klog.Errorf("Create ZoneSettiings %s failed: %v", setting.Short, enconterError)
	return enconterError
}

// @title         GetZoneOfInterface
// @description   temporary add a firewalld port
// @auth      	  author           2021-09-27
// @param         iface    		   string         "e.g. eth0, iface is device name."
// @return        zoneName         string         "Return name (s) of zone the interface is bound to or empty string.."
func (c *DbusClientSerivce) GetZoneOfInterface(iface string) string {
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_GETZONEOFINTERFACE)
	klog.V(4).Infof("Get zone of interface: %v", iface)
	call := obj.Call(apis.ZONE_GETZONEOFINTERFACE, dbus.FlagNoAutoStart, iface)
	return call.Body[0].(string)
}
