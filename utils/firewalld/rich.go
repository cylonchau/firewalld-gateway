package firewalld

import (
	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
)

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
	klog.V(4).Infof("Try to create permanent rich rule %s in zone %s.", rule.ToString(), zone)

	call := obj.Call(apis.CONFIG_ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, rule.ToString())
	if call.Err != nil {
		klog.Errorf("Create permanent rich rule failed: %v", call.Err)
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
	klog.Warningf("Query permanent rich rule %s in zone %s failed: %v", rule.ToString(), zone, encounterError)
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
	printPath(apis.PATH, apis.ZONE_QUERYRICHRULE)
	klog.V(4).Infof("Try to query rich rule %s in zone %s.", rule.ToString, zone)

	call := obj.Call(apis.ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())
	if len(call.Body) <= 0 || !call.Body[0].(bool) {
		klog.Warningf("Query rich rule %s in zone %s failed: %v", rule.ToString(), zone, call.Err.Error())
		return false
	}
	return true
}
