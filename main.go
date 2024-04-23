package main

import (
	"flag"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins"
	"log"
	"os"
	"plugin"

	ep "github.com/evacchi/envoy-ext-server/extproc"
)

func parseArgs(args []string) (cfgFile *string, opts *ep.ProcessingOptions, nonFlagArgs []string) {
	rootCmd := flag.NewFlagSet("root", flag.ExitOnError)
	//conn = rootCmd.String("listen", "tcp://:50051", "The connection string.")

	opts = ep.NewDefaultOptions()

	cfgFile = rootCmd.String("c", "ext-server.yaml", "config file")

	rootCmd.BoolVar(&opts.LogStream, "log-stream", false, "log the stream or not.")
	rootCmd.BoolVar(&opts.LogPhases, "log-phases", false, "log the phases or not.")
	rootCmd.BoolVar(&opts.UpdateExtProcHeader, "update-extproc-header", false, "update the extProc header or not.")
	rootCmd.BoolVar(&opts.UpdateDurationHeader, "update-duration-header", false, "update the duration header or not.")

	rootCmd.Parse(args)
	nonFlagArgs = rootCmd.Args()
	return
}

func main() {
	// cmd subCmd arg, arg2,...
	args := os.Args
	conn := "tcp://:50051"
	if len(args) >= 2 {
		//log.Fatal("Passing a processor is required.")
		conn = args[1]
	}

	cfgFile, opts, nonFlagArgs := parseArgs(os.Args[2:])
	ps := loadFilterChain(cfgFile)
	proc := plugins.NewFilterChain(ps)
	if err := proc.Init(opts, nonFlagArgs, pluginapi.FilterConfig{}); err != nil {
		log.Fatalf("Initialize the processor is failed: %v.", err.Error())
	}
	defer proc.Finish()

	ep.Serve(conn, proc)
}

func loadFilterChain(cfgFile *string) []pluginapi.Plugin {
	config, err := pluginapi.ReadConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	var ps []pluginapi.Plugin
	for _, fc := range config.FilterChains {
		for _, f := range fc.Filters {
			switch f.Type {
			case "built-in":
			case "go-plugin":
				p, err := plugin.Open(f.Path)
				if err != nil {
					log.Fatal(err)
				}
				factory, err := p.Lookup("New")
				if err != nil {
					log.Fatal(err)
				}
				pf := factory.(func() pluginapi.Plugin)
				pluginapi.Register(f.Name, pf)
			}
			p := pluginapi.Instantiate(f)
			ps = append(ps, p)
		}
	}
	return ps
}
