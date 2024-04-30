package extproc

import (
	"fmt"
	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/evacchi/envoy-ext-server/pluginapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"reflect"
)

type ExternalProcessorServer struct {
	name   string
	plugin pluginapi.Plugin
}

func (s *ExternalProcessorServer) Process(srv extprocv3.ExternalProcessor_ProcessServer) error {

	var req *extprocv3.ProcessingRequest
	var resp *extprocv3.ProcessingResponse
	var reqc *reqContext
	var respc *respContext
	var err error
	ctx := srv.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
		}

		req, err = srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}
		switch req := req.Request.(type) {
		case *extprocv3.ProcessingRequest_RequestHeaders:
			resp, reqc, err = s.processRequest_RequestHeaders(req)

		case *extprocv3.ProcessingRequest_RequestBody:
			resp, reqc, err = s.processRequest_RequestBody(reqc, req)

		case *extprocv3.ProcessingRequest_RequestTrailers:
			resp, reqc, err = s.processRequest_RequestTrailers(reqc, req)

		case *extprocv3.ProcessingRequest_ResponseHeaders:
			resp, respc, err = s.processRequest_ResponseHeaders(req)

		case *extprocv3.ProcessingRequest_ResponseBody:
			resp, respc, err = s.processRequest_ResponseBody(respc, req)

		case *extprocv3.ProcessingRequest_ResponseTrailers:
			resp, respc, err = s.processRequest_ResponseTrailers(respc, req)

		default:
			err = fmt.Errorf("unknown request type %v", reflect.TypeOf(req))
		}

		if resp != nil {
			if err := srv.Send(resp); err != nil {
				log.Printf("Send error %v", err)
			}
		}
		if err != nil {
			log.Printf("Phase processing error %v", err)
		}

	} // end for over stream messages
}

// Request

func (s *ExternalProcessorServer) processRequest_RequestHeaders(req *extprocv3.ProcessingRequest_RequestHeaders) (*extprocv3.ProcessingResponse, *reqContext, error) {
	rc, err := newReqContext(req.RequestHeaders.Headers)
	if err != nil {
		return nil, rc, err
	}
	err = s.plugin.OnRequestHeaders(rc)
	if err != nil {
		return nil, rc, err
	}
	return rc.toResponseRequestHeaders(), rc, nil
}

func (s *ExternalProcessorServer) processRequest_RequestBody(rc *reqContext, req *extprocv3.ProcessingRequest_RequestBody) (*extprocv3.ProcessingResponse, *reqContext, error) {
	rc.body.body = req.RequestBody.Body
	err := s.plugin.OnRequestBody(rc)
	if err != nil {
		return nil, rc, err
	}
	return rc.toResponseRequestBody(), rc, nil
}

func (s *ExternalProcessorServer) processRequest_RequestTrailers(rc *reqContext, req *extprocv3.ProcessingRequest_RequestTrailers) (*extprocv3.ProcessingResponse, *reqContext, error) {
	var err error
	rc.trailers, err = toTrailers(req.RequestTrailers.Trailers)
	if err != nil {
		return nil, rc, err
	}
	err = s.plugin.OnRequestTrailers(rc)
	if err != nil {
		return nil, rc, err
	}
	return rc.toResponseRequestTrailers(), rc, nil
}

// Response

func (s *ExternalProcessorServer) processRequest_ResponseHeaders(req *extprocv3.ProcessingRequest_ResponseHeaders) (*extprocv3.ProcessingResponse, *respContext, error) {
	rc, err := newRespContext(req.ResponseHeaders.Headers)
	if err != nil {
		return nil, rc, err
	}
	err = s.plugin.OnResponseHeaders(rc)
	if err != nil {
		return nil, rc, err
	}
	return rc.toResponseResponseHeaders(), rc, nil
}

func (s *ExternalProcessorServer) processRequest_ResponseTrailers(rc *respContext, req *extprocv3.ProcessingRequest_ResponseTrailers) (*extprocv3.ProcessingResponse, *respContext, error) {
	var err error
	rc.trailers, err = toTrailers(req.ResponseTrailers.Trailers)
	if err != nil {
		return nil, rc, err
	}
	err = s.plugin.OnResponseTrailers(rc)
	if err != nil {
		return nil, rc, err
	}
	return rc.toResponseResponseTrailers(), rc, nil
}

func (s *ExternalProcessorServer) processRequest_ResponseBody(rc *respContext, req *extprocv3.ProcessingRequest_ResponseBody) (*extprocv3.ProcessingResponse, *respContext, error) {
	rc.body.body = req.ResponseBody.Body
	err := s.plugin.OnResponseBody(rc)
	if err != nil {
		return nil, rc, nil
	}
	return rc.toResponseResponseBody(), rc, nil
}
