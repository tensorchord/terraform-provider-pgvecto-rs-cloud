package client

import (
	"encoding/json"
	"strings"
	"testing"
)

var errJson = []byte(`
{
	"http_status_code":40005,
	"message":"your api key does not have permission to access this resource"
}
`)

func TestError_UnmarshalJSON(t *testing.T) {
	var e Error
	err := json.Unmarshal(errJson, &e)
	if err != nil {
		t.Errorf("Error.UnmarshalJSON() error = %v", err)
	}
	wantCode := 40005
	if e.HTTPStatusCode != wantCode {
		t.Errorf("Error.UnmarshalJSON() = %v, want %v", e.HTTPStatusCode, wantCode)
	}

	wantMessage := "api key does not have permission"
	if !strings.Contains(e.Message, wantMessage) {
		t.Errorf("Error.UnmarshalJSON() = %v, want %v", e.Message, wantMessage)
	}
}
