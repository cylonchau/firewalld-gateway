package firewalld

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cylonchau/firewalldGateway/log"
)

func (c *DbusClientSerivce) checkZoneName(name string) error {
	if len(name) > 17 {
		log.Error("zone_name is limited to 17 chars:", name)
		return errors.New("zone_name is limited to 17 chars.")
	}
	return nil
}

func splitPortProtocol(portProtocol string) (port, protocol string) {
	if strings.Contains(portProtocol, "/") {
		slices := strings.Split(portProtocol, "/")
		return slices[0], slices[1]
	}
	return portProtocol, "tcp"
}

func checkPort(portProtocol string) (err error) {
	if strings.Contains(portProtocol, "/") {
		slices := strings.Split(portProtocol, "/")
		if strings.Contains(slices[0], "-") {
			portRange := strings.Split(portProtocol, "-")
			frist, _ := strconv.Atoi(portRange[0])
			last, _ := strconv.Atoi(portRange[1])
			if !(frist > 0 && frist <= 65535) && !(last > 0 && last <= 65535) {
				log.Error(fmt.Sprintf("%s,%s", "port range error,expects 1-65535", portRange, "given."))
				return errors.New(fmt.Sprintf("%s,%s", "port range error,expects 1-65535", portRange, "given."))
			}
		}
	} else {
		if strings.Contains(portProtocol, "-") {
			portRange := strings.Split(portProtocol, "-")
			frist, _ := strconv.Atoi(portRange[0])
			last, _ := strconv.Atoi(portRange[1])
			if !(frist > 0 && frist <= 65535) && !(last > 0 && last <= 65535) {
				log.Error(fmt.Sprintf("%s,%s", "port range error,expects 1-65535", portRange, "given."))
				return errors.New(fmt.Sprintf("%s,%s", "port range error,expects 1-65535", portRange, "given."))
			}
		}
	}
	return
}
