package client

import (
	"fmt"
)

type Error struct {
	HTTPStatusCode int `json:"http_status_code,omitempty"`

	// Human-readable message.
	Message string `json:"message,omitempty"`
	Request string `json:"request,omitempty"`

	// Logical operation and nested error.
	Op  string `json:"op,omitempty"`
	Err error  `json:"error,omitempty"`
}

func (err Error) Error() string {
	return fmt.Sprintf("code:%d,Message:%s", err.HTTPStatusCode, err.Message)
}

func (err Error) Is(target error) bool {
	t, ok := target.(Error)
	if !ok {
		return false
	}
	return t.HTTPStatusCode == err.HTTPStatusCode
}
