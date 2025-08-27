package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompter provides generic interactive input functionality
type Prompter struct {
	scanner *bufio.Scanner
}

// NewPrompter creates a new Prompter instance
func NewPrompter() *Prompter {
	return &Prompter{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

type SingleChoiceConfig struct {
	Messages   []string
	Options    []PromptOption
	DefaultKey string
}

type SingleInputConfig struct {
	Message      string
	DefaultValue int
	Validator    func(string) (int, error)
}

type MultipleInputConfig struct {
	HeaderMessages []string
	ParseFunc      func(string) (any, error)
	DoneKeyword    string
	Formatter      func(any) string
}

type PromptOption struct {
	Key   string
	Label string
	Value any
}

// PromptSingleChoice prompts user to select from predefined options
func (p *Prompter) PromptSingleChoice(config SingleChoiceConfig) any {
	for i, msg := range config.Messages {
		if i == len(config.Messages)-1 {
			fmt.Print(msg)
		} else {
			fmt.Println(msg)
		}
	}
	if !p.scanner.Scan() {
		for _, opt := range config.Options {
			if opt.Key == config.DefaultKey {
				return opt.Value
			}
		}
		return nil
	}

	input := strings.ToLower(strings.TrimSpace(p.scanner.Text()))
	fmt.Println()

	if input == "" {
		for _, opt := range config.Options {
			if opt.Key == config.DefaultKey {
				return opt.Value
			}
		}
	}

	for _, opt := range config.Options {
		if opt.Key == input || opt.Label == input {
			return opt.Value
		}
	}

	for _, opt := range config.Options {
		if opt.Key == config.DefaultKey {
			return opt.Value
		}
	}
	return nil
}

// PromptSingleInput prompts user for a single input value with validation
func (p *Prompter) PromptSingleInput(config SingleInputConfig) int {
	fmt.Print(config.Message)
	if !p.scanner.Scan() {
		return config.DefaultValue
	}

	input := strings.TrimSpace(p.scanner.Text())
	if input == "" {
		return config.DefaultValue
	}

	if value, err := config.Validator(input); err == nil {
		return value
	}

	fmt.Printf("Invalid input, using default %d\n", config.DefaultValue)
	return config.DefaultValue
}

// PromptMultipleInput prompts user for multiple inputs until done keyword is entered
func (p *Prompter) PromptMultipleInput(config MultipleInputConfig) []any {
	var results []any

	for _, msg := range config.HeaderMessages {
		fmt.Println(msg)
	}

	for {
		fmt.Print("> ")
		if !p.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(p.scanner.Text())
		if input == config.DoneKeyword {
			fmt.Println()
			break
		}

		result, err := config.ParseFunc(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		results = append(results, result)
		fmt.Printf("Added: %s\n", config.Formatter(result))
	}

	return results
}
