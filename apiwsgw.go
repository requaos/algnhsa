package algnhsa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	pe "github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
)

var (
	errAPIGatewayWebsocketUnexpectedRequest = errors.New("expected APIGatewayWebsocketProxyRequest event")
)

func newAPIGatewayWebsocketRequest(ctx context.Context, payload []byte, opts *Options) (lambdaRequest, error) {
	var event events.APIGatewayWebsocketProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return lambdaRequest{}, err
	}
	if event.RequestContext.AccountID == "" {
		msg := map[string]interface{}{}
		_ = json.Unmarshal(payload, &msg)
		return lambdaRequest{}, pe.WithMessage(errAPIGatewayWebsocketUnexpectedRequest, fmt.Sprintf("%+v", msg["requestContext"]))
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

	if opts.UseProxyPath {
		req.Path = path.Join("/", event.PathParameters["proxy"])
	}

	return req, nil
}
