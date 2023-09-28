package slogger

import (
	"os"

	"github.com/abiosoft/ishell"
)

// ReadLine print tip and wait input.
func ReadLine(c ishell.Actions, val string) string {
	c.Print(val)
	return c.ReadLine()
}

// Info print text and wrap line.
func Info(c ishell.Actions, msg string) {
	c.Printf("%s\n", msg)
}

// Infof printf prints to output using string format.
func Infof(c *ishell.Context, format string, vals ...any) {
	c.Printf(format, vals...)
}

// Warn printf warning message like "Warning: ..."
func Warn(c ishell.Actions, err error) {
	c.Printf("Warning: %s\n", err)
}

// Die printf error message and exit.
func Die(c ishell.Actions, err error) {
	c.Printf("Error: %s\n", err)
	os.Exit(1)
}

// Done printf done message and exit.
func Done(c ishell.Actions, msg string) {
	c.Printf("Done: %s\n", msg)
	os.Exit(0)
}
