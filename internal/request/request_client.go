package request

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mvcris/biziu/internal/parser"
)

type RequestClient struct {
	Properties *parser.Properties
}

type ResponseData struct {
	Status int
	Body   []byte
	Time   int64
}

func NewRequestClient(properties *parser.Properties) *RequestClient {
	return &RequestClient{
		Properties: properties,
	}
}

func (r *RequestClient) DoRequest() *ResponseData {
	postBody, err := json.Marshal(r.Properties.Body)

	if err != nil {
		//@TODO: invalid json, should be return a 500 with invalid json message error?
		panic(err)
	}

	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(r.Properties.Method, r.Properties.Url, responseBody)

	if err != nil {
		//@TODO: invalid properties passed to RequestClient,
		panic(err)
	}

	defer req.Body.Close()
	for key := range r.Properties.Header {
		req.Header.Set(key, r.Properties.Header[key])
	}

	client := &http.Client{}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return &ResponseData{
			Status: 500,
			Body:   nil,
			Time:   time.Since(start).Milliseconds(),
		}
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return &ResponseData{
			Status: resp.StatusCode,
			Body:   nil,
			Time:   time.Since(start).Milliseconds(),
		}
	}

	duration := time.Since(start).Milliseconds()

	return &ResponseData{
		Status: resp.StatusCode,
		Time:   duration,
		Body:   body,
	}
}
