package algnhsa

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var (
	Log                                     *log.Logger
	errAPIGatewayWebsocketUnexpectedRequest = errors.New("expected APIGatewayWebsocketProxyRequest event")
)

func init() {
	Log = log.New(os.Stdout, "Socket Request:", log.LstdFlags)
}

func newAPIGatewayWebsocketRequest(ctx context.Context, payload []byte, opts *Options) (lambdaRequest, error) {
	var event events.APIGatewayWebsocketProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return lambdaRequest{}, err
	}
	if event.RequestContext.APIID == "" || event.RequestContext.EventType == "" {
		return lambdaRequest{}, errAPIGatewayWebsocketUnexpectedRequest
	}
	Log.Printf("Event Details: %s %s %s", event.RequestContext.EventType, event.RequestContext.RouteKey, event.RequestContext.Status)

	var overriddenPath bool
	if opts != nil {
		if v, ok := opts.actionPathOverrideMap[strings.ToLower(event.RequestContext.EventType)]; ok {
			event.Path = v.Path
			event.HTTPMethod = v.HTTPMethod
			overriddenPath = true
		}
	}

	req := lambdaRequest{
		HTTPMethod:                      event.HTTPMethod,
		Path:                            event.Path,
		QueryStringParameters:           event.QueryStringParameters,
		MultiValueQueryStringParameters: event.MultiValueQueryStringParameters,
		Headers:                         event.Headers,
		MultiValueHeaders:               event.MultiValueHeaders,
		Body:                            event.Body,
		IsBase64Encoded:                 event.IsBase64Encoded,
		SourceIP:                        event.RequestContext.Identity.SourceIP,
		Context:                         newWebsocketProxyRequestContext(ctx, event),
	}
	Log.Printf("Event Request Details: %+v", req)

	if opts.UseProxyPath && !overriddenPath {
		req.Path = path.Join("/", event.PathParameters["proxy"])
	}

	return req, nil
}
