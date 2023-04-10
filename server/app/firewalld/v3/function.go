package v3

import (
	"context"
	"time"

	"github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
)

func batchFunction(c context.Context) {
	defer c.Done()
	b := c.Value("action_obj")
	delayTime := c.Value("delay_time").(uint32)
	eventName := c.Value("event_name").(string)
	tName := batch_processor.RandName()
	switch b.(type) {

	case apis.ZoneDst:
		obj := b.(apis.ZoneDst)
		var event batch_processor.Event
		switch eventName {
		case batch_processor.ENABLE_MASQUERADE,
			batch_processor.DISABLE_MASQUERADE,
			batch_processor.SET_DEFAULT_ZONE:
			event = batch_processor.Event{
				EventName: eventName,
				Host:      obj.Host,
				TaskName:  tName,
				Task:      obj.Zone,
			}
		}

		if delayTime > 0 {
			batch_processor.P.AddAfter(tName, time.Duration(delayTime)*time.Second, event)
		} else {
			batch_processor.P.Add(tName, event)
		}
	case apis.PortQuery:
		port := b.(apis.PortQuery)
		delayTime := b.(int)
		eventName := b.(string)
		event := batch_processor.Event{
			EventName: eventName,
			Host:      port.Ip,
			TaskName:  tName,
			Task:      port,
		}
		if delayTime > 0 {
			batch_processor.P.AddAfter(tName, time.Duration(delayTime)*time.Second, event)
		} else {
			batch_processor.P.Add(tName, event)
		}
	case apis.ServiceQuery:
		service := b.(apis.ServiceQuery)
		event := batch_processor.Event{
			EventName: eventName,
			Host:      service.Ip,
			TaskName:  tName,
			Task:      service,
		}
		if delayTime > 0 {
			batch_processor.P.AddAfter(tName, time.Duration(delayTime)*time.Second, event)
		} else {
			batch_processor.P.Add(tName, event)
		}
	case apis.BatchSettingQuery:
		var event batch_processor.Event
		switch eventName {
		case batch_processor.RELOAD_FIREWALD:
			obj := b.(string)
			event = batch_processor.Event{
				EventName: eventName,
				Host:      obj,
				TaskName:  tName,
			}
		}

		if delayTime > 0 {
			batch_processor.P.AddAfter(tName, time.Duration(delayTime)*time.Second, event)
		} else {
			batch_processor.P.Add(tName, event)
		}
	}
}
