package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/gstones/moke-kit/demo/internal/client"
)

var options struct {
	host    string
	tcpHost string
}

const (
	DefaultHost    = "localhost:8081"
	DefaultTcpHost = "localhost:8888"
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
		"grpc http service (<host>:<port>)",
	)

	rootCmd.PersistentFlags().StringVar(
		&options.tcpHost,
		"tcp_host",
		DefaultTcpHost,
		"zinx service (<host>:<port>)",
	)

	sGrpc := &cobra.Command{
		Use:   "grpc",
		Short: "Run an interactive grpc client",
		Run: func(cmd *cobra.Command, args []string) {
			client.RunGrpc(options.host)
		},
	}
	sZinx := &cobra.Command{
		Use:   "zinx",
		Short: "Run an interactive zinx client",
		Run: func(cmd *cobra.Command, args []string) {
			client.RunZinx(options.tcpHost)
		},
	}
	rootCmd.AddCommand(sGrpc, sZinx)
	_ = rootCmd.ExecuteContext(context.Background())
}
