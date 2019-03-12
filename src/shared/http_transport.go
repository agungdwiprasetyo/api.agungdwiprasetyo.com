package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/agungdwiprasetyo/go-utils"
	"github.com/agungdwiprasetyo/go-utils/debug"
)

type Transport struct {
}

type monitoring struct {
	UserAgent string
	Method    string
	Path      string
	Duration  int64
}

func (t *Transport) RoundTrip(request *http.Request) (*http.Response, error) {
	start := time.Now()
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		me := utils.NewMultiError()
		me.Append("server", fmt.Errorf("Server is not reachable. Server not working. Try again later"))

		respBody := httpResponse{
			Code:    http.StatusBadGateway,
			Message: "Bad Gateway",
			Errors:  me.ToMap(),
		}

		bg, _ := json.Marshal(respBody)
		response = &http.Response{
			Status:        "200 OK",
			StatusCode:    http.StatusBadGateway,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBuffer(bg)),
			ContentLength: int64(len(bg)),
			Request:       request,
			Header:        http.Header{"Content-Type": []string{"application/json"}}, //make(http.Header, 0),
		}
	}

	elapsed := time.Since(start)

	var monit = monitoring{
		UserAgent: request.UserAgent(),
		Method:    request.Method,
		Path:      request.URL.Path,
		Duration:  elapsed.Nanoseconds(),
	}

	debug.PrintJSON(monit)

	return response, nil
}