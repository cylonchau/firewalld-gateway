package firewalld

import (
	"github.com/godbus/dbus/v5"

	"github.com/cylonchau/firewalld-gateway/apis"
)

/************************************************** rich rule area ***********************************************************/

// @title         GetRichRules
// @description   Get list of rich-language rules in zone.
// @auther      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @return        zoneName         string         "Returns name of zone to which the interface was bound."
// @return        error            error          "Possible errors: INVALID_ZONE"
func (c *DbusClientSerivce) GetRichRules(zone string) (ruleList []*apis.Rule, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_GETRICHRULES)
	call := obj.Call(apis.ZONE_GETRICHRULES, dbus.FlagNoAutoStart, zone)
	c.eventLogFormat.encounterError = call.Err

	if c.eventLogFormat.encounterError == nil {
		list, ok := call.Body[0].([]string)
		if ok {
			for _, value := range list {
				ruleList = append(ruleList, apis.StringToRule(value))
			}
			c.eventLogFormat.Format = ListResourceSuccessFormat
			c.eventLogFormat.resource = ruleList
			c.printResourceEventLog()
			return
		}
	}
	return nil, c.eventLogFormat.encounterError
}

// @title         AddRichRule
// @description   temporary Add rich language rule into zone.
// @auther      	  author           2021-09-29
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @param         timeout    	   int	          "Timeout, if timeout is non-zero, the operation will be active only for the amount of seconds."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) AddRichRule(zone string, rule *apis.Rule, timeout uint32) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	c.printPath(apis.ZONE_ADDRICHRULE)
	call := obj.Call(apis.ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString(), timeout)

	c.eventLogFormat.encounterError = call.Err
	if c.eventLogFormat.encounterError != nil {
		c.eventLogFormat.Format = CreateResourceFailedFormat
		c.printResourceEventLog()
		return c.eventLogFormat.encounterError
	}
	c.eventLogFormat.Format = CreateResourceSuccessFormat
	c.printResourceEventLog()
	return nil
}

// @title         PermanentAddRichRule
// @description   Permanently Add rich language rule into zone.
// @auther      	  author           2021-10-05
// @param         zone    	       sting 		  "If zone is empty string, use default zone. e.g. public|dmz..  ""
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) AddPermanentRichRule(zone string, rule *apis.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = CreatePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	path, err := c.generatePath(zone, apis.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printPath(apis.CONFIG_ZONE_ADDRICHRULE)

		call := obj.Call(apis.CONFIG_ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, rule.ToString())

		if c.eventLogFormat.encounterError = call.Err; c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = CreatePermanentResourceSuccessFormat
			c.printResourceEventLog()
			return c.eventLogFormat.encounterError
		}
	}
	c.eventLogFormat.Format = CreatePermanentResourceFailedFormat
	c.printResourceEventLog()
	return nil
}

// @title         RemoveRichRule
// @description   temporary Remove rich rule from zone.
// @auther      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: INVALID_ZONE, INVALID_RULE, NOT_ENABLED, INVALID_COMMAND"
func (c *DbusClientSerivce) RemoveRichRule(zone string, rule *apis.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = RemoveResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()

	obj := c.client.Object(apis.INTERFACE, apis.PATH)

	c.printPath(apis.ZONE_REOMVERICHRULE)
	call := obj.Call(apis.ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())

	if c.eventLogFormat.encounterError = call.Err; c.eventLogFormat.encounterError != nil {
		c.eventLogFormat.Format = RemoveResourceFailedFormat
		c.printResourceEventLog()
		return c.eventLogFormat.encounterError
	}
	c.eventLogFormat.Format = RemoveResourceSuccessFormat
	c.printResourceEventLog()
	return nil
}

// @title         PermanentAddRichRule
// @description   Permanently Add rich language rule into zone.
// @auther      	  author           2021-10-05
// @param         zone    	       sting 		  "If zone is empty string, use default zone. e.g. public|dmz..  ""
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) RemovePermanentRichRule(zone string, rule *apis.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = RemovePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()
	path, err := c.generatePath(zone, apis.ZONE_PATH)

	c.eventLogFormat.encounterError = err
	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printPath(apis.CONFIG_ZONE_REOMVERICHRULE)
		call := obj.Call(apis.CONFIG_ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart, rule.ToString())

		if c.eventLogFormat.encounterError = call.Err; c.eventLogFormat.encounterError == nil {
			c.eventLogFormat.Format = RemovePermanentResourceSuccessFormat
			c.printResourceEventLog()
			return nil
		}
	}
	c.eventLogFormat.Format = RemovePermanentResourceFailedFormat
	c.printResourceEventLog()
	return c.eventLogFormat.encounterError
}

// @title         PermanentQueryRichRule
// @description   Check Permanent Configurtion whether rich rule rule has been added in zone.
// @auther      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        bool             bool           "Possible errors: INVALID_ZONE, INVALID_RULE"
func (c *DbusClientSerivce) QueryPermanentRichRule(zone string, rule *apis.Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = QueryPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	path, err := c.generatePath(zone, apis.ZONE_PATH)

	if c.eventLogFormat.encounterError = err; c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		c.printPath(apis.CONFIG_ZONE_QUERYRICHRULE)
		call := obj.Call(apis.CONFIG_ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, rule.ToString())

		if c.eventLogFormat.encounterError = call.Err; c.eventLogFormat.encounterError == nil && (len(call.Body) == 0 || call.Body[0].(bool)) {
			c.eventLogFormat.Format = QueryPermanentResourceSuccessFormat
			c.printResourceEventLog()
			return true
		}
	}
	c.eventLogFormat.Format = QueryPermanentResourceFailedFormat
	c.printResourceEventLog()
	return false
}

// @title         QueryRichRule
// @description   Check whether rich rule is already has.
// @auther      	  author           2021-10-05
// @param         zone    		   string         "If zone is empty string, use default zone. e.g. public|dmz..  "
// @param         rule    	   	   rule	          "rule, rule is rule struct."
// @return        bool             bool           "Possible errors: INVALID_ZONE, INVALID_RULE"
func (c *DbusClientSerivce) QueryRichRule(zone string, rule *apis.Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	c.printPath(apis.ZONE_QUERYRICHRULE)
	call := obj.Call(apis.ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())

	if c.eventLogFormat.encounterError = call.Err; c.eventLogFormat.encounterError == nil && call.Body[0].(bool) {
		c.eventLogFormat.Format = QueryResourceSuccessFormat
		c.printResourceEventLog()
		return true
	}

	c.eventLogFormat.Format = QueryResourceFailedFormat
	c.printResourceEventLog()
	return false
}
