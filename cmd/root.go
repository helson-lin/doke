package cmd

import (
	"fmt"
	"os"

	"github.com/helson-lin/doke/i18n"
	"github.com/spf13/cobra"
)

var (
	version string = "0.0.2"

	rootCmd = &cobra.Command{
		Use:   "doke",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			// 如果没有提供子命令，显示帮助信息
			cmd.Help()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	// 在执行前设置国际化文本
	updateRootCmdText()
	updateAllCommandsText()
	return rootCmd.Execute()
}

// 更新根命令的国际化文本
func updateRootCmdText() {
	rootCmd.Short = i18n.T("root.short")
	rootCmd.Long = i18n.T("root.long")
}

// 更新所有命令的国际化文本
func updateAllCommandsText() {
	for _, cmd := range rootCmd.Commands() {
		switch cmd.Name() {
		case "clear":
			cmd.Short = i18n.T("clear.short")
			cmd.Long = i18n.T("clear.long")
		case "inspect":
			cmd.Short = i18n.T("inspect.short")
			cmd.Long = i18n.T("inspect.long")
		case "proxy":
			cmd.Short = i18n.T("proxy.short")
			cmd.Long = i18n.T("proxy.long")
		case "version":
			cmd.Short = i18n.T("version.short")
			cmd.Long = i18n.T("version.long")
		case "lang":
			cmd.Short = i18n.T("lang.short")
			cmd.Long = i18n.T("lang.long")
		case "command":
			cmd.Short = i18n.T("command.short")
			cmd.Long = i18n.T("command.long")
		case "completion":
			cmd.Short = i18n.T("completion.short")
			cmd.Long = i18n.T("completion.long")
		case "help":
			cmd.Short = i18n.T("help.short")
			cmd.Long = i18n.T("help.long")
		}
	}

	// 更新版本标志的描述
	if versionFlag := rootCmd.PersistentFlags().Lookup("version"); versionFlag != nil {
		versionFlag.Usage = i18n.T("version.print")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, i18n.T("version.print"))
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// 在每次命令执行前更新翻译
		updateAllCommandsText()

		// 检查是否使用了版本标志
		if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
			fmt.Printf(i18n.T("version.current", version))
			fmt.Println()
			os.Exit(0)
		}
		return nil
	}
	// rootCmd.PersistentFlags().StringVarP(&containerId, "id", "i", "", "Docker container id")
	rootCmd.MarkFlagRequired("id")
}
