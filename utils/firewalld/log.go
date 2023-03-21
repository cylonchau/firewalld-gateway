package firewalld

import (
	"reflect"

	"k8s.io/klog/v2"
)

type logFormat struct {
	Format         string
	resourceType   string
	encounterError error
	resource       interface{}
}

const (
	CallRemoteFormat  = "Call %s D-Bus: %s"
	NotFount          = "Resource %s not fount"
	PermanentNotFount = "Permament resource %s not fount"

	// create
	CreateResourceStartFormat            = "Trying to create %s on %s: %v"
	CreateResourceSuccessFormat          = "Create %s on %s successed, %v"
	CreateResourceFailedFormat           = "Create %s on %s failed: %v"
	CreatePermanentResourceStartFormat   = "Trying to create permanent %s on %s: %v"
	CreatePermanentResourceSuccessFormat = "Create permanent %s on %s successed, %v"
	CreatePermanentResourceFailedFormat  = "Create permanent %s on %s failed: %v"

	// query
	QueryResourceStartFormat            = "Trying to query %s on %s: %v"
	QueryResourceSuccessFormat          = "Query %s on %s: %v"
	QueryResourceFailedFormat           = "Query %s on %s failed: %v"
	QueryPermanentResourceStartFormat   = "Trying to query permanent %s on %s: %v"
	QueryPermanentResourceSuccessFormat = "Query permanent %s on %s: %v"
	QueryPermanentResourceFailedFormat  = "Query permanent %s on %s failed: %v"
	QueryNotFount                       = "Resource %s not fount"

	// list
	ListResourceStartFormat            = "Trying list available %s on %s"
	ListResourceSuccessFormat          = "List of %s on %s: %v"
	ListResourceFailedFormat           = "List available %s on %s failed: %v"
	ListPermanentResourceStartFormat   = "Trying list permanent %s on %s"
	ListPermanentResourceSuccessFormat = "List permanent %s on %s: %v"
	ListPermanentResourceFailedFormat  = "List permanent %s on %s failed: %v"

	// delete
	RemoveResourceStartFormat            = "Trying remove %s on %s: %v"
	RemoveResourceSuccessFormat          = "Remove %s on %s successed: %v"
	RemoveResourceFailedFormat           = "Remove %s on %s failed: %v"
	RemovePermanentResourceStartFormat   = "Trying remove permanent %s on %s: %v"
	RemovePermanentResourceSuccessFormat = "Remove permanent %s on %s successed: %v"
	RemovePermanentResourceFailedFormat  = "Remove permanent %s on %s failed: %v"

	// switch
	SwitchResourceStartFormat            = "Trying %s masquerade on %s"
	SwitchPermanentResourceStartFormat   = "Trying permament %s masquerade on %s"
	SwitchResourceSuccessFormat          = "Operation %v masquerade on %s successed"
	SwitchPermanentResourceSuccessFormat = "Operation %v permanent masquerade on %s successed"
	SwitchResourceFailedFormat           = "Operation %v masquerade on %s failed: %v"
	SwitchPermanentResourceFailedFormat  = "Operation %v permanent masquerade on %s failed: %v"

	// zone
	ZoneDefaultStartFormat   = "Trying set default zone on %s to %s "
	ZoneDefaultSuccessFormat = "Set default zone on %s to %s successed"
	ZoneDefaultFailedFormat  = "Set default zone on %s to %s failed: %v"
)

func (c *DbusClientSerivce) printPath(interfaceName string) {
	klog.V(5).Infof(CallRemoteFormat, c.ip, interfaceName)
}

func (c *DbusClientSerivce) printResourceEventLog() {
	if reflect.DeepEqual(c.eventLogFormat, logFormat{}) {
		klog.Errorf("Log format is nil")
	} else {
		switch c.eventLogFormat.Format {
		case NotFount,
			PermanentNotFount:
			klog.Warningf(c.eventLogFormat.Format, c.eventLogFormat.resource)
		case QueryPermanentResourceFailedFormat,
			QueryResourceFailedFormat:
			klog.Warningf(c.eventLogFormat.Format,
				c.eventLogFormat.resourceType, c.ip, c.eventLogFormat.encounterError)
		case RemoveResourceFailedFormat,
			RemovePermanentResourceFailedFormat,
			ListResourceFailedFormat,
			ListPermanentResourceFailedFormat,
			CreatePermanentResourceFailedFormat,
			CreateResourceFailedFormat,
			ZoneDefaultFailedFormat:
			klog.Errorf(c.eventLogFormat.Format,
				c.eventLogFormat.resourceType, c.ip, c.eventLogFormat.encounterError)
		case SwitchResourceFailedFormat,
			SwitchPermanentResourceFailedFormat:
			klog.Warningf(c.eventLogFormat.Format,
				c.eventLogFormat.resourceType, c.ip, c.eventLogFormat.encounterError)
		case ListResourceStartFormat,
			ListPermanentResourceStartFormat,
			SwitchResourceStartFormat, SwitchPermanentResourceStartFormat,
			SwitchResourceSuccessFormat, SwitchPermanentResourceSuccessFormat:
			klog.V(4).Infof(c.eventLogFormat.Format, c.eventLogFormat.resourceType, c.ip)
		case ZoneDefaultStartFormat,
			ZoneDefaultSuccessFormat:
			klog.V(4).Infof(c.eventLogFormat.Format, c.ip, c.eventLogFormat.resource)
		default:
			klog.V(4).Infof(c.eventLogFormat.Format,
				c.eventLogFormat.resourceType, c.ip, c.eventLogFormat.resource)
		}
	}
}
