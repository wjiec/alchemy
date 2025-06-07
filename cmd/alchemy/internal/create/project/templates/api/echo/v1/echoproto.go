package v1

import (
	"github.com/spf13/pflag"

	"github.com/wjiec/alchemy/cmd/alchemy/internal/template"
)

type EchoProto struct{}

func (e *EchoProto) Init(*pflag.FlagSet, template.FlagConstrain) error {
	return nil
}

func (e *EchoProto) Path() string { return "api/echo/v1/echo.proto" }
func (e *EchoProto) Body() string { return echoProtoTemplate }

const echoProtoTemplate = `syntax = "proto3";

package echo.v1;
option go_package = "echo/v1";

import "google/api/annotations.proto";


service EchoService {
  rpc Echo(EchoRequest) returns (EchoResponse) {
    option (google.api.http) = {
      get: '/api/v1/echo'
    };
  }
}

message EchoRequest {
  string text = 1;
}

message EchoResponse {
  string text = 1;
}
`
