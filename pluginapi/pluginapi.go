package pluginapi

import (
	ep "github.com/evacchi/envoy-ext-server/extproc"
)

type Plugin interface {
	Init(opts *ep.ProcessingOptions, nonFlagArgs []string) error
	Finish()

	ep.RequestProcessor
}
