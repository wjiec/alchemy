package alchemy

import (
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// JsonDecoder implements request bodies JSON decoding for HTTP requests.
type JsonDecoder struct {
	runtime.JSONPb
}

// Decoder returns a function that decodes JSON data from an HTTP request body into the provided value.
//
// It handles field selection based on the route description's RequestField.
func (j *JsonDecoder) Decoder(req *http.Request) func(any) error {
	desc, _ := RouteDescFromContext(req.Context())
	return func(raw any) error {
		if req.Header.Get("Content-Type") != j.ContentType(raw) {
			return nil
		}

		if len(desc.RequestField.Name) != 0 {
			raw = desc.RequestField.Accessor(raw)
		}

		if err := j.NewDecoder(req.Body).Decode(raw); err != nil && err != io.EOF {
			return status.Errorf(codes.InvalidArgument, "%v", err)
		}
		return nil
	}
}

// JsonEncoder implements response encoding to JSON format for HTTP responses.
type JsonEncoder struct {
	runtime.JSONPb
}

// Encoder returns a function that encodes a value as JSON in an HTTP response.
func (j *JsonEncoder) Encoder(w http.ResponseWriter, _ *http.Request) func(any) ([]byte, error) {
	return func(resp any) ([]byte, error) {
		w.Header().Set("Content-Type", j.ContentType(resp))
		return j.Marshal(resp)
	}
}
