package request

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/mvcris/biziu/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestHttpCLientSuccess(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	})

	go http.ListenAndServe(":3005", nil)
	parser := parser.NewParser("../../example.json")

	httpClient := NewRequestClient(&parser.Content.Properties)
	res := httpClient.DoRequest()
	assert.Equal(t, string(res.Body), "hello")
	assert.Equal(t, res.Status, 200)
	assert.NotPanics(t, func() { httpClient.DoRequest() }, "not panic")
}

func TestHttpClientPanic(t *testing.T) {
	t.Run("invalid json file", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("function should panic")
			}
		}()
		parser := parser.NewParser("/tmp/invalid_template.json")
		httpClient := NewRequestClient(&parser.Content.Properties)
		httpClient.DoRequest()
		defer os.Remove("/tmp/invalid_template.json")
	})

}
