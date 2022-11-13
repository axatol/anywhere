package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type QueryParameters map[string]string

func (q *QueryParameters) String() string {
	result := make([]string, len(*q))
	index := 0
	for name, value := range *q {
		result[index] = fmt.Sprintf(url.QueryEscape(name), url.QueryEscape(value))
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

func Request[T any](url string, config RequestConfig) (*T, error) {
	if config.Method == "" {
		config.Method = http.MethodGet
	}

	if config.Body == nil {
		config.Body = http.NoBody
	}

	parameters := QueryParameters(config.Parameters)

	request, err := http.NewRequest(config.Method, url+parameters.String(), config.Body)
	if config.Headers != nil {
		for key, value := range config.Headers {
			request.Header.Add(key, value)
		}
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
