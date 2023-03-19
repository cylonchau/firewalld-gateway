package app

import (
	"fmt"
	"io/ioutil"

	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/config"
	"github.com/cylonchau/firewalldGateway/server/app/router"

	"github.com/gin-gonic/gin"
)

var http *gin.Engine

func init() {
	gin.DefaultWriter = ioutil.Discard
	gin.DisableConsoleColor()
}

func NewAPIController() (err error) {
	http = gin.New()
	router.RegisteredRouter(http)
	klog.V(2).Infof("Listening and serving HTTP on %s:%s", config.CONFIG.Address, config.CONFIG.Port)
	if err = http.Run(fmt.Sprintf("%s:%s", config.CONFIG.Address, config.CONFIG.Port)); err != nil {
		return err
	}
	return
}
