package kodik

import (
	"fmt"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
)

// ColorizedOutput provides colorized output functions
type ColorizedOutput struct {
	enabled bool
}

// NewColorizedOutput creates a new colorized output instance
func NewColorizedOutput(enabled bool) *ColorizedOutput {
	return &ColorizedOutput{enabled: enabled}
}

func (c *ColorizedOutput) colorize(color, text string) string {
	if !c.enabled {
		return text
	}
	return color + text + Reset
}

// Success prints a success message in green
func (c *ColorizedOutput) Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(c.colorize(Green, "✓ "+msg))
}

// Error prints an error message in red
func (c *ColorizedOutput) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(c.colorize(Red, "✗ "+msg))
}

// Warning prints a warning message in yellow
func (c *ColorizedOutput) Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(c.colorize(Yellow, "⚠ "+msg))
}

// Info prints an info message in blue
func (c *ColorizedOutput) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(c.colorize(Blue, "ℹ "+msg))
}

// Header prints a header message in bold
func (c *ColorizedOutput) Header(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(c.colorize(Bold+Cyan, "=== "+msg+" ==="))
}

// Summary prints a summary box
func (c *ColorizedOutput) Summary(title string, items []string) {
	fmt.Println(c.colorize(Bold+Cyan, "=== "+title+" ==="))
	for _, item := range items {
		fmt.Println(c.colorize(Green, "  • "+item))
	}
	fmt.Println()
}

// Global colorized output instance
var Output = NewColorizedOutput(true)
