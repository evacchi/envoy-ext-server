package main

import (
	"github.com/evacchi/envoy-ext-server/extproc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	epb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Serve(port int, processor extproc.RequestProcessor) {
	if processor == nil {
		log.Fatalf("cannot process request stream without `processor`")
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(1000)}
	s := grpc.NewServer(sopts...)

	name := processor.GetName()
	opts := processor.GetOptions()
	extprocSvr := extproc.NewGenericExtProcServer(
		name,
		processor,
		opts,
	)
	epb.RegisterExternalProcessorServer(s, extprocSvr)
	hpb.RegisterHealthServer(s, &extproc.HealthServer{})

	log.Printf("Starting ExtProc(%s) on port %d\n", name, port)

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
