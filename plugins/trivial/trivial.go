package trivial

import (
	"encoding/json"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"log"
)

type trivial struct {
	pluginapi.DefaultPlugin
}

func New(config pluginapi.FilterConfig) pluginapi.Plugin {
	return &trivial{}
}

func (t *trivial) OnRequestHeaders(req pluginapi.RequestContext) error {
	doProcess := req.Headers().Get("x-trivial-process")
	if string(doProcess.Value) == "yes" {
		req.Headers().SetRaw("x-custom-header", []byte("hello from the trivial plugin"))
	}
	return nil
}

func (t *trivial) OnRequestBody(req pluginapi.RequestContext) error {
	kv := map[string]string{}
	err := json.Unmarshal(req.Body().Get(), &kv)
	if err != nil {
		log.Fatal(err)
	}
	kv["message"] = "hello from trivial plugin"
	r, err := json.Marshal(kv)
	if err != nil {
		log.Fatal(err)
	}
	req.Body().Set(r)
	return nil
}
