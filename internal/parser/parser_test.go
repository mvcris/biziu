package parser

// Basic imports
import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserImportConfigFile(t *testing.T) {
	parser := NewParser("../../example.json")
	assert.Equal(t, parser.Content.Type, "http")
	assert.Equal(t, parser.Content.Properties.Method, "GET")
	assert.Equal(t, parser.Content.Properties.Url, "http://localhost:3005")
}

func TestParserImportInvalidFIle(t *testing.T) {
	t.Run("file not exists", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("function should panic")
			}
		}()
		NewParser("invalid.json")
	})

	t.Run("invalid json file", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("function should panic")
			}
		}()
		f, _ := os.Create("/tmp/invalid_json_file.json")
		defer f.Close()
		f.Write([]byte("invalid jso content"))
		NewParser("/tmp/invalid_json_file.json")
		defer os.Remove("/tmp/invalid_json_file.json")
	})
}
