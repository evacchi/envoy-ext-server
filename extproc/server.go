package extproc

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"

	epb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Serve(listen string, processor RequestProcessor) {
	if processor == nil {
		log.Fatalf("cannot process request stream without `processor`")
	}

	conn := strings.Split(listen, ":")
	if len(conn) != 2 {
		log.Fatalf("invalid listen address: %s", listen)
	}
	lis, err := net.Listen(conn[0], conn[1])
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(1000)}
	s := grpc.NewServer(sopts...)

	name := processor.GetName()
	opts := processor.GetOptions() // TODO: figure out command line overrides
	extproc := &GenericExtProcServer{
		name:      name,
		processor: processor,
		options:   opts,
	}
	epb.RegisterExternalProcessorServer(s, extproc)
	hpb.RegisterHealthServer(s, &HealthServer{})

	log.Printf("Starting ExtProc(%s) on %s\n", name, listen)

	go s.Serve(lis)

	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	sig := <-gracefulStop
	log.Printf("caught sig: %+v", sig)
	log.Println("Wait for 1 second to finish processing")
	lis.Close()

	time.Sleep(1 * time.Second)
}
