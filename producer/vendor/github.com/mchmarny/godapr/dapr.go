package client

import (
	"time"
)

// StateData represents simplified dapr state item
type StateData struct {
	Key      string            `json:"key"`
	Value    interface{}       `json:"value"`
	Etag     string            `json:"etag,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Options  *StateOptions     `json:"options,omitempty"`
}

// StateOptions is the dapr state data option for StateData
type StateOptions struct {
	Concurrency string       `json:"concurrency,omitempty"`
	Consistency string       `json:"consistency,omitempty"`
	RetryPolicy *RetryPolicy `json:"retryPolicy,omitempty"`
}

// RetryPolicy holds the StateOptions retry policy
type RetryPolicy struct {
	Threshold int32         `json:"threshold,omitempty"`
	Pattern   string        `json:"pattern,omitempty"`
	Interval  time.Duration `json:"interval,omitempty"`
}
