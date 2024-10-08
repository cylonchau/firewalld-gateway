package server

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/server/app"
	"github.com/cylonchau/firewalld-gateway/utils/migration"
	model2 "github.com/cylonchau/firewalld-gateway/utils/model"
)

type Options struct {
	ConfigFile string
	AppName    string
	h          bool
	migration  bool
	sqlDriver  string
	errCh      chan error
}

func NewOptions() *Options {
	return &Options{}
}

// NewProxyCommand creates a *cobra.Command object with default parameters
func NewProxyCommand() *cobra.Command {
	opts := NewOptions()

	cmd := &cobra.Command{
		Use: "",
		Long: `The firewalld-gateway is a central controller as firewallds. 
run only host, docker, kubernetes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			PrintFlags(cmd.Flags())
			if err := opts.Complete(); err != nil {
				return fmt.Errorf("failed complete: %w", err)
			}

			if err := opts.Run(); err != nil {
				klog.ErrorS(err, "Error running "+config.CONFIG.AppName)
				return err
			}

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	fs.AddGoFlagSet(flag.CommandLine) // for --boot-id-file and --machine-id-file

	_ = cmd.MarkFlagFilename("config", "yaml", "yml", "json")

	return cmd
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ConfigFile, "config", "./firewalld-gateway.toml", "The path to the configuration file.")
	fs.BoolVar(&o.migration, "migration", false, "Inital database and tables.")
	fs.StringVar(&o.sqlDriver, "sql-driver", "", "enable which sql backend.")

}

func PrintFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}

func (o *Options) Complete() error {
	if len(o.ConfigFile) == 0 {
		klog.InfoS("Warning, all flags other than --config, --write-config-to, and --cleanup are deprecated, please begin using a config file ASAP")
	}
	// Load the config file here in Complete, so that Validate validates the fully-resolved config.
	if len(o.ConfigFile) > 0 {
		err := config.InitConfiguration(o.ConfigFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Options) Run() error {
	if o.migration {
		return migration.Migration(o.sqlDriver)
	}

	if !config.CONFIG.MySQL.IsEmpty() || !config.CONFIG.SQLite.IsEmpty() {
		if err := model2.InitDB(config.CONFIG.DatabaseDriver); err != nil {
			return err
		}
	}

	return app.NewHTTPSever()
}
