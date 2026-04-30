package agent

import (
	"fmt"
	"io"

	"github.com/cyradin/govnocode/internal/llm"
	"github.com/fatih/color"
)

var (
	defaultColor     = color.New(color.FgBlue)
	taskColor        = color.New(color.FgWhite)
	llmThinkingColor = color.New(color.FgYellow)
	llmMessageColor  = color.New(color.FgWhite)
)

func printDelimiter(out io.Writer) error {
	if _, err := defaultColor.Fprint(out, "-------------------------------------------------\n"); err != nil {
		return err
	}

	return nil
}

func printHeader(out io.Writer, header string) error {
	if _, err := defaultColor.Fprint(out, header, "\n\n"); err != nil {
		return err
	}

	return nil
}

func printTask(out io.Writer, task string) error {
	if err := printDelimiter(out); err != nil {
		return err
	}

	if err := printHeader(out, "Task:"); err != nil {
		return err
	}

	if _, err := taskColor.Fprint(out, task, "\n"); err != nil {
		return err
	}

	return nil
}

func printLLMResponse(out io.Writer, results <-chan llm.ChatResult) error {
	var thinkingMode = false

	if err := printDelimiter(out); err != nil {
		return err
	}

	if err := printHeader(out, "LLM Response:"); err != nil {
		return err
	}

	for result := range results {
		if err := result.Err; err != nil {
			return fmt.Errorf("send message to llm :%w", err)
		}

		msg := result.Resp.Message

		if msg.Thinking != "" {
			if !thinkingMode {
				if _, err := llmThinkingColor.Fprint(out, "Thinking...\n\n"); err != nil {
					return err
				}

				thinkingMode = true
			}

			if _, err := llmThinkingColor.Fprint(out, msg.Thinking); err != nil {
				return err
			}
		}

		if msg.Content != "" {
			if thinkingMode {
				if _, err := llmThinkingColor.Fprint(out, "\n\n"); err != nil {
					return err
				}

				thinkingMode = false
			}

			if _, err := llmMessageColor.Fprint(out, msg.Content); err != nil {
				return err
			}
		}
	}

	return nil
}
