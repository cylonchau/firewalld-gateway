package v3

import (
	"context"
	"time"

	"github.com/cylonchau/firewalld-gateway/server/batch_processor"
	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

func batchFunction(c context.Context) {
	defer c.Done()
	b := c.Value("action_obj")
	delayTime := c.Value("delay_time").(uint32)
	eventName := c.Value("event_name").(string)
	tName := batch_processor.RandName()
	switch b.(type) {

	case query.ZoneDst:
		obj := b.(query.ZoneDst)
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
	case query.PortQuery:
		portRule := b.(query.PortQuery)
		event := batch_processor.Event{
			EventName: eventName,
			Host:      portRule.Ip,
			TaskName:  tName,
			Task:      portRule,
		}
		if delayTime > 0 {
			batch_processor.P.AddAfter(tName, time.Duration(delayTime)*time.Second, event)
		} else {
			batch_processor.P.Add(tName, event)
		}
	case query.ForwardQuery:
		forwardRule := b.(query.ForwardQuery)
		event := batch_processor.Event{
			EventName: eventName,
			Host:      forwardRule.Ip,
			TaskName:  tName,
			Task:      forwardRule,
		}
		if delayTime > 0 {
			batch_processor.P.AddAfter(tName, time.Duration(delayTime)*time.Second, event)
		} else {
			batch_processor.P.Add(tName, event)
		}
	case query.RichQuery:
		richRule := b.(query.RichQuery)
		event := batch_processor.Event{
			EventName: eventName,
			Host:      richRule.Ip,
			TaskName:  tName,
			Task:      richRule,
		}
		if delayTime > 0 {
			batch_processor.P.AddAfter(tName, time.Duration(delayTime)*time.Second, event)
		} else {
			batch_processor.P.Add(tName, event)
		}
	case query.ServiceQuery:
		service := b.(query.ServiceQuery)
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
	case query.BatchSettingQuery:
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
