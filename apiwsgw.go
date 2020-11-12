package algnhsa

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"

	"github.com/aws/aws-lambda-go/events"
)

var (
	Log *log.Logger
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
	if event.RequestContext.APIID == "" {
		return lambdaRequest{}, errAPIGatewayWebsocketUnexpectedRequest
	}
	Log.Printf("Event Details: %s %s %s", event.HTTPMethod, event.Path, event.Body)

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

	if opts.UseProxyPath {
		req.Path = path.Join("/", event.PathParameters["proxy"])
	}

	return req, nil
}
