package llm

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type ToolCall struct {
	Tool string          `json:"tool"`
	Args json.RawMessage `json:"args"`
}

type Parser struct {
	msg string
}

func NewParser(msg string) *Parser {
	return &Parser{
		msg: msg,
	}
}

func (p *Parser) GetTool() (ToolCall, error) {
	if tool, err := p.tryParse(p.msg); err == nil {
		return tool, nil
	}

	code := p.extractCodeBlock(p.msg)
	if code == "" {
		return ToolCall{}, fmt.Errorf("no valid json or code block found")
	}

	return p.tryParse(code)
}

func (p *Parser) tryParse(input string) (ToolCall, error) {
	var call ToolCall

	input = strings.TrimSpace(input)

	if err := json.Unmarshal([]byte(input), &call); err != nil {
		return ToolCall{}, err
	}

	if call.Tool == "" {
		return ToolCall{}, errors.New("missing tool field")
	}

	return call, nil
}

var codeBlockRe = regexp.MustCompile("(?s)```(.*?)```")

func (p *Parser) extractCodeBlock(s string) string {
	m := codeBlockRe.FindStringSubmatch(s)
	if len(m) < 2 { //nolint:mnd
		return ""
	}

	code := strings.TrimSpace(m[1])
	code = strings.TrimPrefix(code, "json")

	return code
}
