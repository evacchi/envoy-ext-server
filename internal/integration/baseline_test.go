package integration

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
)

func BenchmarkBaseline(b *testing.B) {
	cwd, err := os.Getwd()
	println(cwd)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	envoyC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "envoyproxy/envoy:v1.29-latest",
			ExposedPorts: []string{"8000/tcp", "9901/tcp"},
			WaitingFor:   wait.ForListeningPort("8000/tcp"),
			LogConsumerCfg: &testcontainers.LogConsumerConfig{
				Consumers: []testcontainers.LogConsumer{&StdoutLogConsumer{}},
			},
			//Files: []testcontainers.ContainerFile{{
			//	HostFilePath:      cwd + "/testdata/envoy.yaml",
			//	ContainerFilePath: "/etc/envoy/envoy.yaml"}},
			HostConfigModifier: func(config *container.HostConfig) {
				config.Mounts = append(config.Mounts,
					mount.Mount{
						Source: cwd + "/testdata/envoy.yaml",
						Target: "/testdata/envoy.yaml",
						Type:   mount.TypeBind})
			},
		},
		Started: true,
	})

	if err != nil {
		log.Fatalf("Could not start envoy: %s", err)
	}

	println("READY TO RUMBLE")

	defer func() {
		if err := envoyC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	}()

}

// StdoutLogConsumer is a LogConsumer that prints the log to stdout
type StdoutLogConsumer struct{}

// Accept prints the log to stdout
func (lc *StdoutLogConsumer) Accept(l testcontainers.Log) {
	fmt.Print(string(l.Content))
}
