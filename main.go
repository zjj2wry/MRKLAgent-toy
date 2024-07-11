package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

type TerminalTool struct {
	CallbacksHandler callbacks.Handler
}

var _ tools.Tool = TerminalTool{}

// New creates a new terminal tool to execute commands.
func New(opts ...Option) (*TerminalTool, error) {
	options := &options{}

	for _, opt := range opts {
		opt(options)
	}

	return &TerminalTool{}, nil
}

func (t TerminalTool) Name() string {
	return "Terminal"
}

func (t TerminalTool) Description() string {
	return `
	"You have the capability to directly execute terminal commands by terminal tool."
	"Useful for executing system commands directly from the agent. "
	"Action Input required a valid terminal command."`
}

func (t TerminalTool) Call(ctx context.Context, input string) (string, error) {
	if t.CallbacksHandler != nil {
		t.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// Execute the command
	cmd := exec.CommandContext(ctx, "sh", "-c", input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if t.CallbacksHandler != nil {
			t.CallbacksHandler.HandleToolError(ctx, err)
		}
		return "", errors.New(stderr.String())
	}

	result := out.String()

	if t.CallbacksHandler != nil {
		t.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return strings.TrimSpace(result), nil
}

// Option is a function that configures the TerminalTool.
type Option func(*options)

// options is a struct that holds the configuration for TerminalTool.
type options struct{}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	llm, err := openai.New()
	if err != nil {
		return err
	}
	terminal, err := New()
	if err != nil {
		return err
	}
	executor, err := agents.Initialize(
		llm,
		[]tools.Tool{terminal},
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(3),
		agents.WithReturnIntermediateSteps(),
	)
	if err != nil {
		return err
	}
	question := `Write a Golang hello world program in the current directory of my computer, and then execute it.
	You can not write programs on my computer, considering using tools.`
	answer, err := chains.Run(context.Background(), executor, question)
	fmt.Println(answer)
	return err
}
