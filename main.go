package main

import (
	"context"
	"flag"
	ep "github.com/evacchi/envoy-ext-server/extproc"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func parseArgs(args []string) (cfgFile *string) {
	rootCmd := flag.NewFlagSet("root", flag.ExitOnError)
	cfgFile = rootCmd.String("c", "ext-server.yaml", "config file")
	rootCmd.Parse(args)
	return
}

func main() {
	cfgFile := parseArgs(os.Args[1:])
	config, err := pluginapi.ReadConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	proc, err := plugins.NewFilterChainFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	listen := config.Listen
	if listen == "" {
		listen = "tcp://:50051"
	}

	ctx, cancel := context.WithCancel(context.Background())

	ep.Serve(ctx, listen, proc)

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	log.Printf("caught sig: %+v", sig)
	cancel()
}
