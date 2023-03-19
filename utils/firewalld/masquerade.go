package firewalld

import (
	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
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
func (c *DbusClientSerivce) EnableMasquerade(zone string, timeout int) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_ADDMASQUERADE)
	klog.V(4).Infof("Trying enable network masquerade in zone %s .", zone)
	call := obj.Call(apis.ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart, zone, timeout)

	if call.Err != nil && len(call.Body) <= 0 {
		klog.Errorf("Enable network masquerade failed: %v", call.Err.Error())
		return call.Err
	}
	return nil
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
	if path, err := c.generatePath(zone, apis.ZONE_PATH); err == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_ADDMASQUERADE)
		klog.V(4).Infof("Try to permanent enable network masquerade in zone %s .", zone)
		call := obj.Call(apis.CONFIG_ZONE_ADDMASQUERADE, dbus.FlagNoAutoStart)
		err = call.Err
		if err == nil {
			return nil
		}
	}
	klog.Errorf("Permanent enable network masquerade failed: %v", err)
	return err
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
	printPath(apis.PATH, apis.ZONE_REMOVEMASQUERADE)
	klog.V(4).Infof("Trying to disable network masquerade in zone %s .", zone)

	call := obj.Call(apis.ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart, zone)
	if call.Err != nil && len(call.Body) <= 0 {
		klog.Errorf("Disable network masquerade failed:", call.Err.Error())
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

	if path, err := c.generatePath(zone, apis.ZONE_PATH); err == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_REMOVEMASQUERADE)
		klog.V(4).Infof("Trying to disable network masquerade in zone %s.", zone)
		call := obj.Call(apis.CONFIG_ZONE_REMOVEMASQUERADE, dbus.FlagNoAutoStart)
		err = call.Err
		if err == nil && len(call.Body) >= 0 {
			return nil
		}
	}
	klog.Errorf("Disable network masquerade failed: %v", err)
	return err
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

	if path, err := c.generatePath(zone, apis.ZONE_PATH); err == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_ZONE_QUERYMASQUERADE)
		klog.V(4).Infof("Trying to query permanent network masquerade in zone %s.", zone)
		call := obj.Call(apis.CONFIG_ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart)
		err = call.Err
		if call.Err == nil {
			if !call.Body[0].(bool) {
				klog.V(4).Infof("network masquerade in zone %s is disabled.", zone)
				return false, nil
			} else {
				return true, nil
			}
		}
	}
	klog.Errorf("query permanent network masquerade is failed:", err)
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
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.ZONE_QUERYMASQUERADE)
	klog.V(4).Infof("Trying to query permanent network masquerade in zone %s .", zone)

	call := obj.Call(apis.ZONE_QUERYMASQUERADE, dbus.FlagNoAutoStart, zone)
	if call.Err == nil {
		if len(call.Body) <= 0 || !call.Body[0].(bool) {
			klog.V(4).Infof("Network masquerade in zone %s is disabled.")
			return false, call.Err
		}
	}
	return true, nil
}
