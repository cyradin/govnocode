package output

import (
	"io"

	"github.com/fatih/color"
)

type LLMMessage struct {
	Content  string
	Thinking string
}

var (
	defaultColor     = color.New(color.FgBlue)
	taskColor        = color.New(color.FgWhite)
	llmThinkingColor = color.New(color.FgYellow)
	llmMessageColor  = color.New(color.FgWhite)
)

type Printer struct {
	out io.Writer

	defaultColor     *color.Color
	taskColor        *color.Color
	llmThinkingColor *color.Color
	llmMessageColor  *color.Color
	errorColor       *color.Color
}

func NewPrinter(out io.Writer) *Printer {
	return &Printer{
		out: out,

		defaultColor:     color.New(color.FgBlue),
		taskColor:        color.New(color.FgWhite),
		llmThinkingColor: color.New(color.FgYellow),
		llmMessageColor:  color.New(color.FgWhite),
		errorColor:       color.New(color.FgRed),
	}
}

func (p *Printer) PrintUserMessage(msg string) error {
	if err := p.printDelimiter(p.out); err != nil {
		return err
	}

	if err := p.printHeader(p.out, "Message:"); err != nil {
		return err
	}

	if _, err := taskColor.Fprint(p.out, msg, "\n"); err != nil {
		return err
	}

	return nil
}

func (p *Printer) PrintError(err error) error {
	if err := p.printDelimiter(p.out); err != nil {
		return err
	}

	if _, err := p.errorColor.Fprint(p.out, "Error: \n\n"); err != nil {
		return err
	}

	if _, err := p.errorColor.Fprint(p.out, err.Error(), "\n"); err != nil {
		return err
	}

	return nil
}

func (p *Printer) PrintLLMResponse(messages <-chan LLMMessage) error {
	if err := p.printDelimiter(p.out); err != nil {
		return err
	}

	if err := p.printHeader(p.out, "LLM Response:"); err != nil {
		return err
	}

	var (
		isThinking = false
		err        error
	)

	for msg := range messages {
		isThinking, err = p.printLLMThinkingBlock(msg.Thinking, isThinking)
		if err != nil {
			return err
		}

		if err := p.printLLMMessage(msg.Content); err != nil {
			return err
		}
	}

	if _, err := p.llmMessageColor.Fprint(p.out, "\n"); err != nil {
		return err
	}

	return nil
}

func (p *Printer) printLLMThinkingBlock(text string, isThinking bool) (bool, error) {
	if isThinking && text == "" {
		if _, err := llmThinkingColor.Fprint(p.out, "\n\n"); err != nil {
			return false, err
		}

		return false, nil
	}

	if !isThinking && text != "" {
		if _, err := llmThinkingColor.Fprint(p.out, "Thinking...\n\n"); err != nil {
			return false, err
		}

		isThinking = true
	}

	if _, err := llmThinkingColor.Fprint(p.out, text); err != nil {
		return false, err
	}

	return isThinking, nil
}

func (p *Printer) printLLMMessage(text string) error {
	if _, err := llmMessageColor.Fprint(p.out, text); err != nil {
		return err
	}

	return nil
}

func (p *Printer) printDelimiter(out io.Writer) error {
	if _, err := defaultColor.Fprint(out, "-------------------------------------------------\n"); err != nil {
		return err
	}

	return nil
}

func (p *Printer) printHeader(out io.Writer, header string) error {
	if _, err := defaultColor.Fprint(out, header, "\n\n"); err != nil {
		return err
	}

	return nil
}
