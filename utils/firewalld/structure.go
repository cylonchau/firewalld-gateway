package firewalld

import (
	"errors"
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/apis"
	"github.com/cylonchau/firewalldGateway/config"
)

var (
	dbusClient     *DbusClientSerivce
	remotelyBusLck sync.Mutex
	PORT           = "55557"
)

type DbusClientSerivce struct {
	client      *dbus.Conn
	defaultZone string
	ip          string
	port        string
}

func NewDbusClientService(addr string) (*DbusClientSerivce, error) {
	remotelyBusLck.Lock()
	defer remotelyBusLck.Unlock()
	var (
		encounterError error
		conn           *dbus.Conn
	)
	if config.CONFIG.Port == "" {
		klog.V(5).Infof("Start connect to D-Bus service: %s:%s", addr, PORT)
	} else {
		PORT = config.CONFIG.Dbus_Port
		klog.V(5).Infof("Start connect to D-Bus service: %s:%s", addr, PORT)
	}

	if dbusClient != nil && dbusClient.client.Connected() {
		return dbusClient, nil
	}
	if conn, encounterError = dbus.Connect("tcp:host="+addr+",port="+PORT, dbus.WithAuth(dbus.AuthAnonymous())); encounterError == nil {
		obj := conn.Object(apis.INTERFACE, apis.PATH)
		call := obj.Call(apis.INTERFACE_GETDEFAULTZONE, dbus.FlagNoAutoStart)
		encounterError = call.Err
		if encounterError == nil {
			return &DbusClientSerivce{
				conn,
				call.Body[0].(string),
				addr,
				PORT,
			}, encounterError
		}
	}
	klog.Errorf("Connect to firewalld client failed: %v", encounterError)
	return nil, encounterError
}

/*
 * @title         destroy
 * @description   off firewalld connection.
 * @auth          author    2021-10-31
 */
func (c *DbusClientSerivce) Destroy() {
	err := c.client.Close()
	if err != nil {
		klog.Errorf("Close D-Bus connection failed, %v", err)
	}
}

/************************************************** fw service area ***********************************************************/

/*
 * @title         Reload
 * @description   temporary Add rich language rule into zone.
 * @auth          author           2021-10-05
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) Reload() error {
	obj := c.client.Object(apis.INTERFACE, apis.PATH)
	printPath(apis.PATH, apis.INTERFACE_RELOAD)
	klog.V(4).Infof("Try to reload firewalld runtime.")
	call := obj.Call(apis.INTERFACE_RELOAD, dbus.FlagNoAutoStart)

	if call.Err != nil {
		klog.Errorf("reload firewalld failed: %v", call.Err.Error())
		return call.Err
	}
	return nil
}

/*
 * @title         flush currently zone zoneSettings to default zoneSettings.
 * @description   temporary Add rich language rule into zone.
 * @auth          author           2021-10-05
 * @return        error            error          "Possible errors:
 *                                                      ALREADY_ENABLED"
 */
func (c *DbusClientSerivce) RuntimeFlush(zone string) (encounterError error) {
	if zone == "" {
		zone = c.GetDefaultZone()
	}

	defaultZoneSetting := apis.Settings{
		Target:      "default",
		Description: "reset by " + config.CONFIG.AppName,
		Short:       "public",
		Interface:   nil,
		Service: []string{
			"ssh",
			"dhcpv6-client",
		},
		Port: []*apis.Port{
			&apis.Port{
				Port:     config.CONFIG.Dbus_Port,
				Protocol: "tcp",
			},
		},
	}

	var path dbus.ObjectPath
	if path, encounterError = c.generatePath(zone, apis.ZONE_PATH); encounterError == nil {
		obj := c.client.Object(apis.INTERFACE, path)
		printPath(path, apis.CONFIG_UPDATE)
		klog.V(4).Infof("Try to flush current active zone (%s).", zone)
		call := obj.Call(apis.CONFIG_UPDATE, dbus.FlagNoAutoStart, defaultZoneSetting)
		encounterError = call.Err
		if encounterError == nil || len(call.Body) <= 0 {
			return nil
		}
	}

	klog.Errorf("Flush current zone (%s) failed: %v", zone, encounterError)
	return encounterError
}

// @title         Reload
// @description   temporary Add rich language rule into zone.
// @auth      	  author           2021-10-05
// @return        error            error          "Possible errors: ALREADY_ENABLED"
func (c *DbusClientSerivce) generatePath(zone, interface_path string) (dbus.ObjectPath, error) {
	zoneid := c.getZoneId(zone)
	if zoneid < 0 {
		klog.Errorf("invalid zone:", zone)
		return "", errors.New("invalid zone " + interface_path + zone)
	}
	p := fmt.Sprintf("%s/%d", interface_path, zoneid)
	klog.V(5).Infof("Dbus PATH: %s", p)
	return dbus.ObjectPath(p), nil
}

func printPath(pathName dbus.ObjectPath, interfaceName string) {
	klog.V(5).Infof("Call remote D-Bus: %s/%s", pathName, interfaceName)
}
