package batch_processor

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/cylonchau/firewalld-gateway/apis"
	"github.com/cylonchau/firewalld-gateway/utils/firewalld"
)

var Store map[string]interface{}
var mu sync.Mutex

type Event struct {
	EventName string
	Host      string
	TaskName  string
	Task      interface{}
	errNum    int
}

func init() {
	Store = make(map[string]interface{}, 1024)
}

func StoreAdd(key string, v interface{}) {
	mu.Lock()
	defer mu.Unlock()
	Store[key] = v
}

func StoreDel(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(Store, key)
}

func RandName() string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(62)]
	}
	return "task-" + string(b) + strconv.Itoa(int(time.Now().Unix()))
}

func (e *Event) processEvent() error {
	var (
		incurredError error
		dbusClient    *firewalld.DbusClientSerivce
	)

	if dbusClient, incurredError = firewalld.NewDbusClientService(e.Host); incurredError != nil {
		return incurredError
	}
	defer dbusClient.Destroy()

	defaultZone := dbusClient.GetDefaultZone()
	switch e.EventName {
	case CREATE_PORT:
		query := e.Task.(apis.PortQuery)
		incurredError = dbusClient.AddPort(&query.Port, defaultZone, query.Timeout)
	case REMOVE_PORT:
		query := e.Task.(apis.PortQuery)
		incurredError = dbusClient.RemovePort(&query.Port, defaultZone)
	case CREATE_RICH:
		query := e.Task.(apis.RichQuery)
		incurredError = dbusClient.AddRichRule(defaultZone, query.Rich, query.Timeout)
	case CREATE_FORWARD:
		query := e.Task.(apis.ForwardQuery)
		incurredError = dbusClient.AddForwardPort(defaultZone, query.Timeout, query.Forward)
	case CREATE_SERVICE:
		query := e.Task.(apis.ServiceQuery)
		incurredError = dbusClient.AddService(query.Zone, query.Service, query.Timeout)
	case ENABLE_MASQUERADE:
		query := e.Task.(string)
		incurredError = dbusClient.EnableMasquerade(query, 0)
	case DISABLE_MASQUERADE:
		query := e.Task.(string)
		incurredError = dbusClient.DisableMasquerade(query)
	case RELOAD_FIREWALD:
		incurredError = dbusClient.Reload()
	case FLUSH_SETTING:
		query := e.Task.(string)
		incurredError = dbusClient.RuntimeFlush(query)
	case SET_DEFAULT_ZONE:
		query := e.Task.(string)
		incurredError = dbusClient.SetDefaultZone(query)
	case REMOVE_PROTOCOL:

	default:
		incurredError = errors.New("unkown event")
	}

	return incurredError
}
