package parser

import (
	"encoding/json"
	"os"
)

type Parser struct {
	filePath string
	Content  Content
}

type Properties struct {
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"headers"`
	Body   map[string]any    `json:"body"`
}

type Options struct {
	Requests    uint32 `json:"requests"`
	Concurrency uint32 `json:"concurrency"`
	Nodes       uint32 `json:"nodes"`
	Port        uint16 `json:"port"`
}

type Content struct {
	Type       string     `json:"type"`
	Properties Properties `json:"properties"`
	Options    Options    `json:"options"`
}

func NewParser(filePath string) *Parser {
	parser := &Parser{
		filePath: filePath,
	}
	parser.parseTemplateFile()
	return parser
}

func (p *Parser) parseTemplateFile() {
	bytes, err := os.ReadFile(p.filePath)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &p.Content)
	if err != nil {
		panic(err)
	}
}
