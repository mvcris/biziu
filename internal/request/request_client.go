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
	Parser parser.Parser
}

type ResponseData struct {
	Status int
	Body   []byte
	Time   time.Duration
}

func NewRequestClient(parser *parser.Parser) *RequestClient {
	return &RequestClient{
		Parser: *parser,
	}
}

func (r *RequestClient) DoRequest() (*ResponseData, error) {
	postBody, err := json.Marshal(r.Parser.Content.Properties.Body)

	if err != nil {
		return nil, err
	}

	responseBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(r.Parser.Content.Properties.Method, r.Parser.Content.Properties.Url, responseBody)

	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	for key := range r.Parser.Content.Properties.Header {
		req.Header.Set(key, r.Parser.Content.Properties.Header[key])
	}

	client := &http.Client{}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)

	response := &ResponseData{
		Status: resp.StatusCode,
		Time:   duration,
		Body:   body,
	}

	return response, nil
}
