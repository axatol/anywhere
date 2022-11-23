package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/axatol/anywhere/config"
)

type Headers map[string]string

type QueryParameters map[string]string

func (q *QueryParameters) String() string {
	result := make([]string, len(*q))
	index := 0
	for name, value := range *q {
		result[index] = fmt.Sprintf("%s=%s", url.QueryEscape(name), url.QueryEscape(value))
		index += 1
	}

	return fmt.Sprintf("?%s", strings.Join(result, "&"))
}

type RequestConfig struct {
	Method     string
	Parameters map[string]string
	Headers    map[string]string
	Body       io.Reader
}

func Request[T any](url string, cfg RequestConfig) (*T, error) {
	if cfg.Method == "" {
		cfg.Method = http.MethodGet
	}

	if cfg.Body == nil {
		cfg.Body = http.NoBody
	}

	parameters := QueryParameters(cfg.Parameters)
	endpoint := url + parameters.String()
	request, err := http.NewRequest(cfg.Method, endpoint, cfg.Body)
	if cfg.Headers != nil {
		for key, value := range cfg.Headers {
			request.Header.Add(key, value)
		}
	}

	response, err := http.DefaultClient.Do(request)
	config.Log.Debugw("response",
		"request_method", request.Method,
		"request_url", request.URL,
		"response_status", response.StatusCode,
		"response_content_length", response.ContentLength,
	)

	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// config.Log.Debugw("response body", "body", string(raw))

	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
