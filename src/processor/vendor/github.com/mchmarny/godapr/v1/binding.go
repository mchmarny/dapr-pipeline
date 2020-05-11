package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// InvokeBindingWithData posts the in content to specified binding
func (c *Client) InvokeBindingWithData(ctx trace.SpanContext, binding string, in []byte) (out []byte, err error) {

	url := fmt.Sprintf("%s/v1.0/bindings/%s", c.url, binding)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(in))
	if err != nil {
		err = errors.Wrapf(err, "error creating invoking binding request: %s", url)
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

// InvokeBinding serializes data and invokes InvokeBindingWithData
func (c *Client) InvokeBinding(ctx trace.SpanContext, binding string, data interface{}) (out []byte, err error) {
	b, _ := json.Marshal(data)
	return c.InvokeBindingWithData(ctx, binding, b)
}
