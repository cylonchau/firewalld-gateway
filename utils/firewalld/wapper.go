package firewalld

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"k8s.io/klog/v2"
)

func splitPortProtocol(portProtocol string) (port, protocol string) {
	if strings.Contains(portProtocol, "/") {
		slices := strings.Split(portProtocol, "/")
		return slices[0], slices[1]
	}
	return portProtocol, "tcp"
}

func checkPort(portProtocol string) (err error) {
	errSting := "port range error,expects 1-65535, given %s."
	if strings.Contains(portProtocol, "/") {
		slices := strings.Split(portProtocol, "/")
		if strings.Contains(slices[0], "-") {
			portRange := strings.Split(portProtocol, "-")
			frist, _ := strconv.Atoi(portRange[0])
			last, _ := strconv.Atoi(portRange[1])
			if !(frist > 0 && frist <= 65535) && !(last > 0 && last <= 65535) {
				klog.Errorf(errSting, portRange)
				return
			}
		}
	} else {
		if strings.Contains(portProtocol, "-") {
			portRange := strings.Split(portProtocol, "-")
			frist, _ := strconv.Atoi(portRange[0])
			last, _ := strconv.Atoi(portRange[1])
			if !(frist > 0 && frist <= 65535) && !(last > 0 && last <= 65535) {
				klog.Errorf(errSting, portRange)
				return errors.New(fmt.Sprintf(errSting, portRange))
			}
		}
	}
	return
}
