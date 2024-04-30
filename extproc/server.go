package extproc

import (
	"context"
	"fmt"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/evacchi/envoy-ext-server/plugins"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"

	epb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
)

func ServeFromConfig(ctx context.Context, cfgFile string) {
	config, err := pluginapi.ReadConfig(cfgFile)
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
	Serve(ctx, listen, proc)
}

func Serve(ctx context.Context, listen string, plugin pluginapi.Plugin) error {
	conn := strings.Split(listen, "://")
	if len(conn) != 2 {
		return fmt.Errorf("invalid listen address: %s", listen)
	}
	lis, err := net.Listen(conn[0], conn[1])
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(1000)}
	s := grpc.NewServer(sopts...)

	name := "test"
	extproc := &ExternalProcessorServer{
		name:   name,
		plugin: plugin,
	}

	epb.RegisterExternalProcessorServer(s, extproc)
	hpb.RegisterHealthServer(s, &HealthServer{})

	log.Printf("Starting ExtProc(%s) on %s\n", name, listen)

	go s.Serve(lis)

	<-ctx.Done()
	lis.Close()
	return nil
}
