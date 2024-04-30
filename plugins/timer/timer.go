package timer

import (
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"strconv"
	"time"
)

type timer struct {
	pluginapi.DefaultPlugin
	started time.Time
}

func New(fc pluginapi.FilterConfig) pluginapi.Plugin {
	return &timer{
		started: time.Now(),
	}
}

func (t *timer) OnRequestHeaders(req pluginapi.RequestContext) error {
	req.Headers().SetRaw("x-extproc-started-ns", []byte(strconv.FormatInt(t.started.UnixNano(), 10)))
	return nil
}

func (t *timer) OnResponseHeaders(resp pluginapi.ResponseContext) error {
	finished := time.Now()
	duration := time.Since(t.started)

	resp.Headers().SetRaw("x-extproc-started-ns", []byte(strconv.FormatInt(t.started.UnixNano(), 10)))
	resp.Headers().SetRaw("x-extproc-finished-ns", []byte(strconv.FormatInt(finished.UnixNano(), 10)))
	resp.Headers().SetRaw("x-upstream-duration-ns", []byte(strconv.FormatInt(duration.Nanoseconds(), 10)))
	return nil
}
