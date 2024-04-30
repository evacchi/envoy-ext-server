package e2e

import (
	"log"
	"os/exec"
	"testing"
)

func setupEnvoy(t *testing.T, cfg string) *exec.Cmd {
	cmd := exec.Command("envoy", "-c", cfg)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return cmd
}
