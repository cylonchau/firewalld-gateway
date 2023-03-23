package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/server"
)

func main() {
	command := server.NewProxyCommand()
	flagset := flag.CommandLine
	klog.InitFlags(flagset)
	pflag.CommandLine.AddGoFlagSet(flagset)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
