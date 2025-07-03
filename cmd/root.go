package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version string = "0.0.2"

	rootCmd = &cobra.Command{
		Use:   "doke",
		Short: "Convert Docker container config to docker run command",
		Run: func(cmd *cobra.Command, args []string) {
			// 如果没有提供子命令，显示帮助信息
			cmd.Help()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print version information")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// 检查是否使用了版本标志
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf("Current Version: %s\n", version)
			os.Exit(0)
		}
	}
	// rootCmd.PersistentFlags().StringVarP(&containerId, "id", "i", "", "Docker container id")
	rootCmd.MarkFlagRequired("id")
}
