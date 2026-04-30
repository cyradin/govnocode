package agent

import (
	"io"

	"github.com/fatih/color"
)

type PrinterLLMMessage struct {
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
}

func NewPrinter(out io.Writer) *Printer {
	return &Printer{
		out: out,

		defaultColor:     color.New(color.FgBlue),
		taskColor:        color.New(color.FgWhite),
		llmThinkingColor: color.New(color.FgYellow),
		llmMessageColor:  color.New(color.FgWhite),
	}
}

func (p *Printer) PrintTaskText(task string) error {
	if err := p.printDelimiter(p.out); err != nil {
		return err
	}

	if err := p.printHeader(p.out, "Task:"); err != nil {
		return err
	}

	if _, err := taskColor.Fprint(p.out, task, "\n"); err != nil {
		return err
	}

	return nil
}

func (p *Printer) PrintLLMResponse(messages <-chan PrinterLLMMessage) error {
	var thinkingMode = false

	if err := p.printDelimiter(p.out); err != nil {
		return err
	}

	if err := p.printHeader(p.out, "LLM Response:"); err != nil {
		return err
	}

	for msg := range messages {
		if msg.Thinking != "" {
			if !thinkingMode {
				if _, err := llmThinkingColor.Fprint(p.out, "Thinking...\n\n"); err != nil {
					return err
				}

				thinkingMode = true
			}

			if _, err := llmThinkingColor.Fprint(p.out, msg.Thinking); err != nil {
				return err
			}
		}

		if msg.Content != "" {
			if thinkingMode {
				if _, err := llmThinkingColor.Fprint(p.out, "\n\n"); err != nil {
					return err
				}

				thinkingMode = false
			}

			if _, err := llmMessageColor.Fprint(p.out, msg.Content); err != nil {
				return err
			}
		}
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
