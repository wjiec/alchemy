package alchemy

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// PathDecoder implements request path parameter decoding for HTTP requests.
//
// It extracts path variables from the URL and converts them to structured data.
type PathDecoder struct{}

// Decoder returns a function that extracts path parameters from an HTTP
// request and decodes them into the provided value.
func (p *PathDecoder) Decoder(req *http.Request) func(any) error {
	desc, _ := RouteDescFromContext(req.Context())
	return func(msg any) error {
		requiredNames := make(map[string]struct{})
		for _, name := range desc.PathParameters {
			requiredNames[name] = struct{}{}
		}

		for k, v := range mux.Vars(req) {
			delete(requiredNames, k)
			if err := runtime.PopulateFieldFromPath(msg.(proto.Message), k, v); err != nil {
				return status.Errorf(codes.InvalidArgument, "type mismatch, path parameter: %s, error: %v", k, err)
			}
		}

		for missingName := range requiredNames {
			return status.Errorf(codes.InvalidArgument, "missing path parameter %s", missingName)
		}
		return nil
	}
}
