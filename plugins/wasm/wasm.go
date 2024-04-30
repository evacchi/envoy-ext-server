package wasm

import (
	"bytes"
	"context"
	"fmt"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"github.com/mitchellh/mapstructure"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"log"
	"os"
)

type Config struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type wasm struct {
	pluginapi.DefaultPlugin
	rt wazero.Runtime
	m  wazero.CompiledModule
}

func New(config pluginapi.FilterConfig) pluginapi.Plugin {
	cfg := &Config{}
	err := mapstructure.Decode(config.Config, cfg)
	if err != nil {
		log.Fatal(err)
	}
	wasmBin, err := os.ReadFile(cfg.Path)
	if err != nil {
		log.Fatal(err)
	}

	c := context.Background()
	rt := wazero.NewRuntime(c)
	wasi_snapshot_preview1.MustInstantiate(c, rt)

	m, err := rt.CompileModule(c, wasmBin)
	if err != nil {
		log.Fatal(err)
	}

	return &wasm{rt: rt, m: m}
}

func (w *wasm) OnRequestBody(req pluginapi.RequestContext) error {
	c := context.Background()

	stdout := bytes.NewBuffer([]byte{})
	headers := req.Headers().All()
	args := make([]string, len(headers)+1)
	args[0] = "wasm"
	for i, h := range headers {
		args[i+1] = fmt.Sprintf("%s=%s", h.Key, h.Value)
	}

	cfg := wazero.NewModuleConfig().
		WithStdout(stdout).
		WithStdin(bytes.NewReader(req.Body().Get())).
		WithArgs(args...)

	_, err := w.rt.InstantiateModule(c, w.m, cfg)
	if err != nil {
		return err
	}

	req.Body().Set(stdout.Bytes())
	return nil
}
