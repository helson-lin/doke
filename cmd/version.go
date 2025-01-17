package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"}, // 添加别名 v
	Short:   "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Current Version: %s", version)
	},
}
