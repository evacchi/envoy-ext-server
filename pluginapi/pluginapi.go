package pluginapi

import ep "github.com/wrossmorrow/envoy-extproc-sdk-go"

type Plugin interface {
	Init(opts *ep.ProcessingOptions, nonFlagArgs []string) error
	Finish()

	ep.RequestProcessor
}
