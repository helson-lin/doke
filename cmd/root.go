package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version string = "0.0.1"

	rootCmd = &cobra.Command{
		Use:   "doke",
		Short: "Convert Docker container config to docker run command",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// rootCmd.PersistentFlags().StringVarP(&containerId, "id", "i", "", "Docker container id")
	rootCmd.MarkFlagRequired("id")
}
