package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalldGateway/server"
)

var (
	configuration string
	h             bool
)

func BuildInitFlags() {
	flagset := flag.CommandLine
	flagset.StringVar(&configuration, "f", "./firewalld-gateway.conf", "set configuration file.")
	flagset.BoolVar(&h, "h", false, "Prints a short help text and exists.")
	flagset.Usage = usage
	klog.InitFlags(flagset)
	flag.Parse()
}

func cmdPrompt(str string) {
	fmt.Fprintf(os.Stderr, "firewall-api: invalid option, %s\n", str)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: firewalld-gateway [-f configfile] [-h help]

Options
`)
	flag.PrintDefaults()
}

func main() {
	command := server.NewProxyCommand()
	flagset := flag.CommandLine
	klog.InitFlags(flagset)
	pflag.CommandLine.AddGoFlagSet(flagset)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
