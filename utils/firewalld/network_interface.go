package firewalld

import (
	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
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
func (c *DbusClientSerivce) BindInterface(zone, interface_name string) (list string, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_ADDINTERFACE)
	klog.V(4).Infof("Trying to bind interface %s to rule in zone %s .", interface_name, zone)

	call := obj.Call(apis.ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)
	if call.Err != nil {
		klog.Errorf("bind interface %s failed: %v", interface_name, call.Err.Error())
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
	if path, err = c.generatePath(zone, apis.ZONE_PATH); err == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_ADDINTERFACE)
		klog.V(4).Infof("Trying to permanent bind interface %s to rule in zone %s.", interface_name, zone)
		call := obj.Call(apis.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interface_name)
		err = call.Err
		if call.Err == nil {
			return nil
		}
	}
	klog.Errorf("Permanent bind interface %s, failed: %v", interface_name, err)
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
func (c *DbusClientSerivce) QueryInterface(zone, interface_name string) (b bool, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_QUERYINTERFACE)
	klog.V(4).Infof("Trying to query interface %s bind in zone %s .", interface_name, zone)
	call := obj.Call(apis.ZONE_QUERYINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)

	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		klog.Errorf("Query interface %s bind failed: %v", interface_name, call.Err)
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
func (c *DbusClientSerivce) PermanentQueryInterface(zone, interface_name string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	path, err := c.generatePath(zone, apis.ZONE_PATH)
	if err == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_ADDINTERFACE)
		klog.V(4).Infof("Trying to query permanent interface %s bind in zone %s .", interface_name, zone)
		call := obj.Call(apis.CONFIG_ZONE_ADDINTERFACE, dbus.FlagNoAutoStart, interface_name)
		err = call.Err
		if err == nil {
			return nil
		}
	}
	klog.Errorf("Query permanent interface %s bind failed: %v", interface_name, err)
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
func (c *DbusClientSerivce) RemoveInterface(zone, interface_name string) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_REMOVEINTERFACE)
	klog.V(4).Infof("Trying to remove interface %s bind in zone %s.", interface_name, zone)
	call := obj.Call(apis.ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, zone, interface_name)
	if call.Err != nil && len(call.Body) <= 0 {
		klog.Errorf("Remove interface %s bind failed: %v", interface_name, call.Err)
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
	printPath(path, apis.CONFIG_ZONE_REMOVEINTERFACE)
	klog.V(4).Infof("Trying to remove permanent interface %s bind in zone %s.", interface_name, zone)
	call := obj.Call(apis.CONFIG_ZONE_REMOVEINTERFACE, dbus.FlagNoAutoStart, interface_name)
	if call.Err != nil {
		klog.Errorf("Remove permanent interface %s bind failed: %v", interface_name, call.Err)
		return call.Err
	}
	return nil
}
