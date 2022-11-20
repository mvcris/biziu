package parser

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Parser struct {
	filePath string
	Content  Content
}

type Content struct {
	Type       string `json:"type"`
	Properties struct {
		Url    string            `json:"url"`
		Method string            `json:"method"`
		Header map[string]string `json:"headers"`
		Body   map[string]any    `json:"body"`
	} `json:"properties"`
}

func NewParser(filePath string) *Parser {
	parser := &Parser{
		filePath: filePath,
	}
	parser.parseTemplateFile()
	return parser
}

func (p *Parser) parseTemplateFile() {
	jsonFile, err := os.Open(p.filePath)

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byteValue, &p.Content)
	if err != nil {
		panic(err)
	}
}
