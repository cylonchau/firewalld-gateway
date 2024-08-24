package firewalld

import (
	"github.com/godbus/dbus/v5"

	api2 "github.com/cylonchau/firewalld-gateway/api"
)

/************************************************** rich rule area ***********************************************************/

func (c *DbusClientSerivce) GetRichRules(zone string) (ruleList []*api2.Rule, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(api2.INTERFACE, api2.PATH)

	c.printPath(api2.ZONE_GETRICHRULES)
	call := obj.Call(api2.ZONE_GETRICHRULES, dbus.FlagNoAutoStart, zone)
	c.eventLogFormat.encounterError = call.Err

	if c.eventLogFormat.encounterError == nil {
		list, ok := call.Body[0].([]string)
		if ok {
			for _, value := range list {
				ruleList = append(ruleList, api2.StringToRule(value))
			}
			c.eventLogFormat.Format = ListResourceSuccessFormat
			c.eventLogFormat.resource = ruleList
			c.printResourceEventLog()
			return
		}
	}
	return nil, c.eventLogFormat.encounterError
}

func (c *DbusClientSerivce) GetPermanentRichRules(zone string) (ruleList []*api2.Rule, err error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = ListResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(api2.INTERFACE, api2.PATH)

	c.printPath(api2.CONFIG_ZONE_GETRICHRULES)
	call := obj.Call(api2.ZONE_GETRICHRULES, dbus.FlagNoAutoStart, zone)
	c.eventLogFormat.encounterError = call.Err

	if c.eventLogFormat.encounterError == nil {
		list, ok := call.Body[0].([]string)
		if ok {
			for _, value := range list {
				ruleList = append(ruleList, api2.StringToRule(value))
			}
			c.eventLogFormat.Format = ListResourceSuccessFormat
			c.eventLogFormat.resource = ruleList
			c.printResourceEventLog()
			return
		}
	}
	return nil, c.eventLogFormat.encounterError
}

func (c *DbusClientSerivce) AddRichRule(zone string, rule *api2.Rule, timeout uint32) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = CreateResourceStartFormat
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(api2.INTERFACE, api2.PATH)
	c.printPath(api2.ZONE_ADDRICHRULE)
	call := obj.Call(api2.ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString(), timeout)

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

func (c *DbusClientSerivce) AddPermanentRichRule(zone string, rule *api2.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = CreatePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	path, err := c.generatePath(zone, api2.ZONE_PATH)
	c.eventLogFormat.encounterError = err

	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(api2.INTERFACE, path)
		c.printPath(api2.CONFIG_ZONE_ADDRICHRULE)

		call := obj.Call(api2.CONFIG_ZONE_ADDRICHRULE, dbus.FlagNoAutoStart, rule.ToString())

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

func (c *DbusClientSerivce) RemoveRichRule(zone string, rule *api2.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = RemoveResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil

	c.printResourceEventLog()

	obj := c.client.Object(api2.INTERFACE, api2.PATH)

	c.printPath(api2.ZONE_REOMVERICHRULE)
	call := obj.Call(api2.ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())

	if c.eventLogFormat.encounterError = call.Err; c.eventLogFormat.encounterError != nil {
		c.eventLogFormat.Format = RemoveResourceFailedFormat
		c.printResourceEventLog()
		return c.eventLogFormat.encounterError
	}
	c.eventLogFormat.Format = RemoveResourceSuccessFormat
	c.printResourceEventLog()
	return nil
}

func (c *DbusClientSerivce) RemovePermanentRichRule(zone string, rule *api2.Rule) error {
	if zone == "" {
		zone = c.GetDefaultZone()
	}
	// print log
	c.eventLogFormat.Format = RemovePermanentResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()
	path, err := c.generatePath(zone, api2.ZONE_PATH)

	c.eventLogFormat.encounterError = err
	if c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(api2.INTERFACE, path)
		c.printPath(api2.CONFIG_ZONE_REOMVERICHRULE)
		call := obj.Call(api2.CONFIG_ZONE_REOMVERICHRULE, dbus.FlagNoAutoStart, rule.ToString())

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

func (c *DbusClientSerivce) QueryPermanentRichRule(zone string, rule *api2.Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = QueryPermanentResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	path, err := c.generatePath(zone, api2.ZONE_PATH)

	if c.eventLogFormat.encounterError = err; c.eventLogFormat.encounterError == nil {
		obj := c.client.Object(api2.INTERFACE, path)
		c.printPath(api2.CONFIG_ZONE_QUERYRICHRULE)
		call := obj.Call(api2.CONFIG_ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, rule.ToString())

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

func (c *DbusClientSerivce) QueryRichRule(zone string, rule *api2.Rule) bool {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	// print log
	c.eventLogFormat.Format = QueryResourceStartFormat
	c.eventLogFormat.resourceType = "rich"
	c.eventLogFormat.resource = rule.ToString()
	c.eventLogFormat.encounterError = nil
	c.printResourceEventLog()

	obj := c.client.Object(api2.INTERFACE, api2.PATH)
	c.printPath(api2.ZONE_QUERYRICHRULE)
	call := obj.Call(api2.ZONE_QUERYRICHRULE, dbus.FlagNoAutoStart, zone, rule.ToString())

	if c.eventLogFormat.encounterError = call.Err; c.eventLogFormat.encounterError == nil && call.Body[0].(bool) {
		c.eventLogFormat.Format = QueryResourceSuccessFormat
		c.printResourceEventLog()
		return true
	}

	c.eventLogFormat.Format = QueryResourceFailedFormat
	c.printResourceEventLog()
	return false
}
