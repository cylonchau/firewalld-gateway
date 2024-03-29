package app

import (
	"fmt"
	"io/ioutil"

	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/server/app/router"
	"github.com/cylonchau/firewalld-gateway/server/batch_processor"

	"github.com/gin-gonic/gin"
)

var http *gin.Engine
var stopCh = make(chan struct{})

func init() {
	gin.DefaultWriter = ioutil.Discard
	gin.DisableConsoleColor()
}

func NewHTTPSever() (err error) {
	http = gin.New()
	router.RegisteredRouter(http)
	klog.V(2).Infof("Listening and serving HTTP on %s:%s", config.CONFIG.Address, config.CONFIG.Port)

	if config.CONFIG.AsyncProcess {
		batch_processor.P = batch_processor.NewProcessor()
		go batch_processor.P.Run()
	}
	if err = http.Run(fmt.Sprintf("%s:%s", config.CONFIG.Address, config.CONFIG.Port)); err != nil {
		return err
	}
	<-stopCh
	return
}
