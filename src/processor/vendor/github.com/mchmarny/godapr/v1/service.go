package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// InvokeServiceWithData invokes the remote service method
func (c *Client) InvokeServiceWithData(ctx trace.SpanContext, service, method string, in []byte) (out []byte, err error) {
	url := fmt.Sprintf("%s/v1.0/invoke/%s/method/%s", c.url, service, method)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(in))
	if err != nil {
		err = errors.Wrapf(err, "error creating invoking request: %s", url)
		return
	}

	content, status, err := c.exec(ctx, req)
	if err != nil {
		err = errors.Wrapf(err, "error executing: %+v", req)
		return
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("invalid response code to %s: %d", url, status)
	}

	return content, nil
}

// InvokeService serializes input data to JSON and invokes InvokeServiceWithData
func (c *Client) InvokeService(ctx trace.SpanContext, service, method string, in interface{}) (out []byte, err error) {
	b, _ := json.Marshal(in)
	return c.InvokeServiceWithData(ctx, service, method, b)
}
