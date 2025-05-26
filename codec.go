package alchemy

import (
	"errors"
	"net/http"
)

// DecoderFactory defines an interface for creating HTTP request decoders.
type DecoderFactory interface {
	// Decoder creates a function that can decode a request data into a given value.
	Decoder(*http.Request) func(any) error
}

// EncoderFactory defines an interface for creating HTTP response encoders.
type EncoderFactory interface {
	// Encoder creates a function that can encode a value into a response
	Encoder(http.ResponseWriter, *http.Request) func(any) ([]byte, error)
}

// CodecFactory defines an interface for encoding and decoding HTTP requests and responses.
type CodecFactory interface {
	DecoderFactory
	EncoderFactory
}

// NewHttpDynamicCodec returns a CodecFactory implementation that can determine the appropriate
// encoding/decoding strategy dynamically based on request properties.
func NewHttpDynamicCodec() CodecFactory {
	return &httpDynamicCodec{
		decoderFactories: []DecoderFactory{
			&JsonDecoder{},
			&MultipartCodec{},
			&QueryDecoder{},
			&PathDecoder{},
		},
		encoderFactories: []EncoderFactory{
			&JsonEncoder{},
		},
	}
}

// httpDynamicCodec implements the CodecFactory interface with dynamic content negotiation capabilities.
type httpDynamicCodec struct {
	decoderFactories []DecoderFactory
	encoderFactories []EncoderFactory
}

// Decoder returns a function that decodes HTTP request data into the provided value.
//
// The decoder function determines the appropriate decoding strategy based on the request.
func (h *httpDynamicCodec) Decoder(req *http.Request) func(any) error {
	return func(v any) error {
		for _, factory := range h.decoderFactories {
			if dec := factory.Decoder(req); dec != nil {
				if err := dec(v); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// Encoder returns a function that encodes a value into an HTTP response.
//
// The encoder function determines the appropriate encoding strategy based on the request and response.
func (h *httpDynamicCodec) Encoder(w http.ResponseWriter, r *http.Request) func(any) ([]byte, error) {
	return func(v any) ([]byte, error) {
		for _, factory := range h.encoderFactories {
			if enc := factory.Encoder(w, r); enc != nil {
				return enc(v)
			}
		}
		return nil, errors.New("no encoder factory found")
	}
}
