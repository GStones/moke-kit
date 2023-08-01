package main

import (
	"context"

	"github.com/spf13/cobra"

	"moke-kit/demo/internal/client"
)

var options struct {
	host string
}

const (
	DefaultHost = "localhost:8081"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "demo",
		Short:   "cli",
		Aliases: []string{"cli"},
	}
	rootCmd.PersistentFlags().StringVar(
		&options.host,
		"host",
		DefaultHost,
		"service (<host>:<port>)",
	)

	shell := &cobra.Command{
		Use:   "shell",
		Short: "Run an interactive auth service client",
		Run: func(cmd *cobra.Command, args []string) {
			client.Shell(options.host)
		},
	}
	rootCmd.AddCommand(shell)
	_ = rootCmd.ExecuteContext(context.Background())
}
