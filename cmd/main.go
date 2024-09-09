package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/server"

	_ "github.com/cylonchau/firewalld-gateway/docs"
)

// @title Uranus API
// @version 1.0
// @description Uranus, distrubed firewall gateway.

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @host localhost:2952
// @BasePath /server
// @schemes http
func main() {
	command := server.NewProxyCommand()
	flagset := flag.CommandLine
	klog.InitFlags(flagset)
	pflag.CommandLine.AddGoFlagSet(flagset)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
