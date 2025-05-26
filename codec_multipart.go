package alchemy

import (
	"net/http"
	"strconv"
	"strings"
	"unsafe"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// MultipartCodec implements request multipart form data decoding for HTTP requests.
type MultipartCodec struct{}

// Decoder returns a function that extracts multipart form data from an HTTP
// request and decodes it into the provided value.
func (m *MultipartCodec) Decoder(req *http.Request) func(any) error {
	desc, _ := RouteDescFromContext(req.Context())
	return func(raw any) error {
		if req.Header.Get("Content-Type") != "multipart/form-data" {
			return nil
		}
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			return status.Errorf(codes.InvalidArgument, "%v", err)
		}

		if len(desc.RequestField.Name) != 0 {
			raw = desc.RequestField.Accessor(raw)
		}

		var filters [][]string
		if len(desc.RequestField.Name) != 0 {
			filters = append(filters, strings.Split(desc.RequestField.Name, "."))
		}
		for _, pathParameter := range desc.PathParameters {
			filters = append(filters, strings.Split(pathParameter, "."))
		}

		values := req.MultipartForm.Value
		for key, files := range req.MultipartForm.File {
			values[key] = make([]string, len(files))
			for i, file := range files {
				values[key][i] = strconv.FormatUint(uint64(uintptr(unsafe.Pointer(file))), 10)
			}
		}

		return runtime.PopulateQueryParameters(raw.(proto.Message), values, utilities.NewDoubleArray(filters))
	}
}
