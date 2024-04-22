package main

import (
	"flag"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins"
	"log"
	"os"

	ep "github.com/evacchi/envoy-ext-server/extproc"
)

var processors = map[string]pluginapi.Plugin{
	"noop":    plugins.NewNoopRequestProcessor(),
	"trivial": plugins.NewTrivialRequestProcessor(),
	"timer":   plugins.NewTimerRequestProcessor(),
	"data":    plugins.NewDataRequestProcessor(),
	"digest":  plugins.NewDigestRequestProcessor(),
	"dedup":   plugins.NewDedupRequestProcessor(),
	"masker":  plugins.NewMaskerRequestProcessor(),
	"echo":    plugins.NewEchoRequestProcessor(),
	"wasm":    plugins.NewWasmRequestProcessor(),
}

func parseArgs(args []string) (conn *string, opts *ep.ProcessingOptions, nonFlagArgs []string) {
	rootCmd := flag.NewFlagSet("root", flag.ExitOnError)
	conn = rootCmd.String("listen", "tcp://:50051", "The connection string.")

	opts = ep.NewDefaultOptions()

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
	if len(args) < 2 {
		log.Fatal("Passing a processor is required.")
	}

	//cmd := args[1]
	//proc, exists := processors[cmd]
	//if !exists {
	//	log.Fatalf("Processor \"%s\" not defined.", cmd)
	//}
	port, opts, nonFlagArgs := parseArgs(os.Args[2:])

	var names []string
	var procs []ep.RequestProcessor
	for n, p := range processors {
		names = append(names, n)
		procs = append(procs, p)
		if err := p.Init(opts, nonFlagArgs); err != nil {
			log.Fatalf("Initialize the processor is failed: %v.", err.Error())
		}
		defer p.Finish()

	}

	proc := plugins.NewMultiplexRequestProcessor(names, procs)
	if err := proc.Init(opts, nonFlagArgs); err != nil {
		log.Fatalf("Initialize the processor is failed: %v.", err.Error())
	}
	defer proc.Finish()

	ep.Serve(*port, proc)
}
