package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// NewBody returns a new *bytes.Reader from the given `content`.
func NewBody(t *testing.T, content interface{}) *bytes.Reader {
	encodedBody, err := json.Marshal(&content)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	return bytes.NewReader(encodedBody)
}

// NewRequest performs a new request over the given http.HandlerFunc, using the given `method` and `body`. Returns the response generated by the handler.
func NewRequest(handlerFunc http.HandlerFunc, method string, url string, body *bytes.Reader) *http.Response {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	handlerFunc(res, req)

	return res.Result()
}
