package main

import (
	"firewall-api/cmd/app"
	"firewall-api/config"
	"firewall-api/log"
	"flag"
	"fmt"
	"os"
)

var (
	configuration string
	h             bool
)

func init() {
	flag.StringVar(&configuration, "f", "./firewalld.conf", "set configuration file.")
	flag.BoolVar(&h, "h", false, "Prints a short help text and exists.")
	flag.Usage = usage
	flag.Parse()
}

func cmdPrompt(str string) {
	fmt.Fprintf(os.Stderr, "firewall-api: invalid option, %s\n", str)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: firewall-api [-f configfile] [-h help]

Options
`)
	flag.PrintDefaults()
}

func main() {
	if h {
		flag.Usage()
		return
	}

	if err := config.InitConfiguration(configuration); err != nil {
		cmdPrompt(err.Error())
		flag.Usage()
		return
	}
	log.New(config.CONFIG.LogLevel)
	app.NewAPIController()
}
