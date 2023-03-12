package server

import (
	"context"
	"encoding/json"
	"fmt"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoyt "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"log"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// AuthorizationServer is the struct definition that acts as the Envoy external authorization service class definition.
//
// When we have values that we need to share across invocations to an instance of the service but not more
// widely, they can be defined here.
type AuthorizationServer struct{}

// denied is a shorthand function for reject a request as unauthorized in some fashion
func denied(code int32, body string) *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &status.Status{Code: code},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoyt.HttpStatus{
					Code: envoyt.StatusCode(code),
				},
				Body: body,
			},
		},
	}
}

// allowed is a shorthand function for generating a success response to Envoy, let this request go through
func allowed() *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: []*core.HeaderValueOption{
					{
						Header: &core.HeaderValue{
							Key:   "x-header-set-by-extauth",
							Value: "bla-bla-bla",
						},
					},
				},
			},
		},
	}
}

// Check implements Envoy Authorization service. Proto file:
// https://github.com/envoyproxy/envoy/blob/main/api/envoy/service/auth/v3/external_auth.proto
func (a *AuthorizationServer) Check(ctx context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {

	// Log the target URL of the request that we are verifying authorization for
	var msgBuilder strings.Builder
	httpRequest := req.Attributes.Request.Http
	msgBuilder.WriteString("/n######### REQUEST #########/n")
	msgBuilder.WriteString("URL:      " + httpRequest.Scheme + ":" + httpRequest.Host + httpRequest.Path + "/n")
	msgBuilder.WriteString("Method:   " + httpRequest.Method + "/n")
	msgBuilder.WriteString("Protocol: " + httpRequest.Protocol + "/n")

	// Log the incoming headers as a formatted JSON structure
	msgBuilder.WriteString("/n=== Headers ===/n")
	jsonBytes, err := json.MarshalIndent(httpRequest.Headers, "", "  ")
	if err == nil {
		msgBuilder.WriteString(string(jsonBytes))
	} else {
		msgBuilder.WriteString("failed to marshal headers: " + err.Error())
	}
	msgBuilder.WriteString("/n===============/n")

	// Log the headers / parameters that were sent just for us and will not be passed through to the
	// "upstream" services deeper inside our private network
	jsonBytes, err = json.MarshalIndent(req.Attributes.ContextExtensions, "", "  ")
	log.Println("/n=== Context Extensions ===/n")
	if err == nil {
		log.Println(string(jsonBytes))
	} else {
		msgBuilder.WriteString("failed to marshal context extensions: " + err.Error())
	}
	fmt.Println("/n==========================/n")

	// Default to allowing the request (for now)
	return allowed(), nil
}
