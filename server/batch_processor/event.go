package batch_processor

import (
	"net"

	"github.com/cylonchau/firewalldGateway/apis"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"
)

type Event struct {
	Name string
	Host net.IP
	Task interface{}
}

func (e *Event) ProcessEvent() error {
	var (
		incurredError error
		dbusClient    *firewalld.DbusClientSerivce
	)
	if dbusClient, incurredError = firewalld.NewDbusClientService(e.Host.String()); incurredError != nil {
		return incurredError
	}
	defer dbusClient.Destroy()
	default_zone := dbusClient.GetDefaultZone()
	switch e.Name {
	case "add-port":
		query := e.Task.(apis.PortQuery)
		incurredError = dbusClient.AddPort(query.Port, default_zone, query.Timeout)
	case "remove-port":
		query := e.Task.(apis.PortQuery)
		incurredError = dbusClient.RemovePort(query.Port, default_zone)
	case "add-richRule":
		query := e.Task.(apis.RichQuery)
		incurredError = dbusClient.AddRichRule(default_zone, query.Rich, query.Timeout)
	case "remove-richRule":
		query := e.Task.(apis.RichQuery)
		incurredError = dbusClient.RemoveRichRule(default_zone, query.Rich)
	case "add-forward":
		query := e.Task.(apis.ForwardQuery)
		incurredError = dbusClient.AddForwardPort(default_zone, query.Timeout, query.Forward)
	case "add-services":
		query := e.Task.(apis.Query)
		incurredError = dbusClient.AddService(default_zone, query.Service, query.Timeout)
	case "remove-protocol":
	case "remove-protocol":
	case "remove-protocol":
	case "remove-protocol":
	case "remove-protocol":
	case "remove-protocol":

	default:
		return nil
	}

	return incurredError
}

func NewEvent(ip net.IP) (*firewalld.DbusClientSerivce, error) {
	var (
		incurredError error
		dbusclient    *firewalld.DbusClientSerivce
	)

	if dbusclient, incurredError = firewalld.NewDbusClientService(ip.String()); incurredError == nil {
		return dbusclient, nil
	}
	return nil, incurredError
}
