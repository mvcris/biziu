package request

import (
	"io"
	"net/http"
	"testing"

	"github.com/mvcris/biziu/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	})

	go http.ListenAndServe(":3000", nil)
	parser := parser.NewParser("../../example.json")

	httpClient := NewRequestClient(parser)
	res, err := httpClient.DoRequest()
	assert.Nil(t, err)
	assert.Equal(t, string(res.Body), "hello")
	assert.Equal(t, res.Status, 200)
	assert.NotPanics(t, func() { httpClient.DoRequest() }, "not panic")
}
