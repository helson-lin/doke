package cmd

import (
	"fmt"

	"github.com/helson-lin/doke/i18n"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"}, // 添加别名 v
	Short:   i18n.T("version.short"),
	Long:    i18n.T("version.long"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(i18n.T("version.current", version))
	},
}
