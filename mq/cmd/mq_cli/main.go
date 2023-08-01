package main

import (
	"github.com/gstones/platform/services/common/mq/pkg/tests"
	"github.com/spf13/cobra"
)

var options struct {
	implementation    string
	mqUrl             string
	topics            []string
	username          string
	appID             string
	deliverySemantics string
	consumerURL       string
	producerUrl       string
}

// Default values for our CLI options
const (
	// Default message queue implementation
	defaultImplementation = "local"

	// Default message queue delivery semantics
	defaultDeliverySemantics = "AtMostOnce"

	// Default url for nsq consumer operations
	defaultConsumerUrl = "127.0.0.1:4161"

	// Default url for nsq producer operations
	defaultProducerUrl = "127.0.0.1:4150"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mq",
		Short: "Message Queue Service CLI",
	}
	{
		// Suite Default Topics
		subTopics := append(options.topics,
			"foo",
			"foo.bar",
			"foo.bar.wiz",
			"foo.bar.wiz.cli",
			"foo.foo",
			"foo.foo.foo",
			"0.1.2.3",
			"0.1.2.3.4.5.6.7.8.9.10",
		)
		pubTopics := append(options.topics,
			"foo",
			"foo.bar",
			"foo.bar.wiz",
			"foo.bar.wiz.cli",
			"foo.foo",
			"foo.foo.foo",
			"0.1.2.3",
			"0.1.2.3.4.5.6.7.8.9.10",
		)
		mqTest := &cobra.Command{
			Use:   "mqTest",
			Short: "Test Suite: Subscribes to, publishes to, then unsubscribes from provided topics.",
			Run: func(cmd *cobra.Command, args []string) {

				// It's necessary to generate our uconfig at this level
				// to retain the ability to use different urls,
				// implementations, or delivery semantics.
				commandConfig := tests.NewMQSuiteConfig(
					options.mqUrl,
					subTopics,
					pubTopics,
					options.implementation,
					options.deliverySemantics,
					options.consumerURL,
					options.producerUrl,
				)
				tests.MQSuite(commandConfig)
			},
		}
		// Command Flags
		{
			mqTest.PersistentFlags().StringVar(
				&options.implementation,
				"impl",
				defaultImplementation,
				"mq implementation - valid choices are kafka, local, nats, or nsq",
			)
			mqTest.PersistentFlags().StringVar(
				&options.mqUrl,
				"url",
				"",
				"URL to connect to",
			)
			mqTest.PersistentFlags().StringVar(
				&options.deliverySemantics,
				"delivery",
				defaultDeliverySemantics,
				"mq delivery semantics - valid choices are AtMostOnce, AtLeastOnce",
			)
			mqTest.PersistentFlags().StringVar(
				&options.consumerURL,
				"consumerUrl",
				defaultConsumerUrl,
				"consumerUrl - url for nsq consumer",
			)
			mqTest.PersistentFlags().StringVar(
				&options.producerUrl,
				"producerUrl",
				defaultProducerUrl,
				"producerUrl - url for nsq producer",
			)
		}
		rootCmd.AddCommand(mqTest)
	}
	rootCmd.Execute()
}
