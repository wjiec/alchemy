package alchemy

import (
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// QueryDecoder implements request query parameter decoding for HTTP requests.
type QueryDecoder struct{}

// Decoder returns a function that extracts query parameters from an HTTP
// request and decodes them into the provided value.
func (q *QueryDecoder) Decoder(req *http.Request) func(any) error {
	desc, _ := RouteDescFromContext(req.Context())
	return func(raw any) error {
		if err := req.ParseForm(); err != nil {
			return status.Errorf(codes.InvalidArgument, "%v", err)
		}

		var filters [][]string
		if len(desc.RequestField.Name) != 0 {
			filters = append(filters, strings.Split(desc.RequestField.Name, "."))
		}
		for _, pathParameter := range desc.PathParameters {
			filters = append(filters, strings.Split(pathParameter, "."))
		}

		return runtime.PopulateQueryParameters(raw.(proto.Message), req.Form, utilities.NewDoubleArray(filters))
	}
}
