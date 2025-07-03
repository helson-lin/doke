package cmd

import (
	"fmt"

	"github.com/helson-lin/doke/i18n"
	"github.com/spf13/cobra"
)

var langCmd = &cobra.Command{
	Use:     "lang [language]",
	Aliases: []string{"language", "locale"},
	Short:   "Set or display current language / 设置或显示当前语言",
	Long: `Set or display the current interface language.
Supported languages: zh (Chinese), en (English)

设置或显示当前界面语言。
支持的语言: zh (中文), en (英文)`,
	Example: `  # Display current language / 显示当前语言
  doke lang
  
  # Set language to English / 设置语言为英文
  doke lang en
  
  # Set language to Chinese / 设置语言为中文
  doke lang zh
  
  # Using aliases / 使用别名
  doke language zh
  doke locale en`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// 显示当前语言和支持的语言列表
			currentLang := i18n.GetCurrentLanguage()
			supportedLangs := i18n.GetSupportedLanguages()

			fmt.Printf("Current language / 当前语言: %s\n", currentLang)
			fmt.Printf("Supported languages / 支持的语言: %v\n", supportedLangs)

			// 显示环境变量设置方法
			fmt.Println("\nTo set language permanently / 永久设置语言:")
			fmt.Println("  export DOKE_LANG=zh    # For Chinese / 中文")
			fmt.Println("  export DOKE_LANG=en    # For English / 英文")
			return
		}

		// 设置语言
		newLang := args[0]
		if err := i18n.SetLanguage(newLang); err != nil {
			fmt.Printf("Error / 错误: %v\n", err)
			fmt.Printf("Supported languages / 支持的语言: %v\n", i18n.GetSupportedLanguages())
			return
		}

		// 根据设置的语言显示成功消息
		if newLang == "zh" {
			fmt.Printf("✅ 语言已设置为中文\n")
			fmt.Println("💡 提示: 这只会影响当前会话。要永久设置，请使用: export DOKE_LANG=zh")
		} else {
			fmt.Printf("✅ Language set to English\n")
			fmt.Println("💡 Tip: This only affects the current session. For permanent setting, use: export DOKE_LANG=en")
		}
	},
}

func init() {
	rootCmd.AddCommand(langCmd)
}
