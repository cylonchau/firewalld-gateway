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
// @version 0.0.9
// @description Uranus, distrubed firewall gateway.

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".

// @host localhost:2952
// @BasePath /
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
