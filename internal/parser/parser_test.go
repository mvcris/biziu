package parser

// Basic imports
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	parser := NewParser("../../example.json")
	assert.Equal(t, parser.Content.Type, "http")
	assert.Equal(t, parser.Content.Properties.Method, "POST")
	assert.Equal(t, parser.Content.Properties.Url, "http://localhost:3000")
}
