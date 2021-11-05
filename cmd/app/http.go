package app

import (
	"fmt"
	"io/ioutil"

	"firewall-api/cmd/app/router"
	"firewall-api/config"
	"firewall-api/log"

	"github.com/gin-gonic/gin"
)

var http *gin.Engine

func init() {
	gin.DefaultWriter = ioutil.Discard
	gin.DisableConsoleColor()
}

func NewAPIController() {
	http = gin.New()
	router.RegisterRouter(http)
	log.Info(fmt.Sprintf("Listening and serving HTTP on %s:%s", config.CONFIG.Address, config.CONFIG.Port))
	http.Run(fmt.Sprintf("%s:%s", config.CONFIG.Address, config.CONFIG.Port))
}
