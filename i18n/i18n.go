package i18n

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	currentLang  = "zh"
	printer      *message.Printer
	translations = make(map[string]map[string]string)
	mu           sync.RWMutex
)

// 支持的语言列表
var supportedLanguages = map[string]language.Tag{
	"zh": language.Chinese,
	"en": language.English,
}

// 初始化国际化系统
func Init() error {
	// 从环境变量获取语言设置
	if lang := os.Getenv("DOKE_LANG"); lang != "" {
		currentLang = lang
	} else if lang := os.Getenv("LANG"); lang != "" {
		// 解析系统语言环境变量
		if strings.HasPrefix(lang, "zh") {
			currentLang = "zh"
		} else {
			currentLang = "en"
		}
	}

	// 加载翻译文件
	if err := loadTranslations(); err != nil {
		return fmt.Errorf("加载翻译文件失败: %v", err)
	}

	// 设置打印器
	if tag, ok := supportedLanguages[currentLang]; ok {
		printer = message.NewPrinter(tag)
	} else {
		printer = message.NewPrinter(language.Chinese)
		currentLang = "zh"
	}

	return nil
}

// 加载翻译文件
func loadTranslations() error {
	// 只加载内置翻译数据
	translations["zh"] = zhTranslations
	translations["en"] = enTranslations
	return nil
}

// 获取翻译文本
func T(key string, args ...interface{}) string {
	mu.RLock()
	defer mu.RUnlock()

	if trans, ok := translations[currentLang]; ok {
		if text, exists := trans[key]; exists {
			if len(args) > 0 {
				return fmt.Sprintf(text, args...)
			}
			return text
		}
	}

	// 如果当前语言没有找到，尝试英文
	if currentLang != "en" {
		if trans, ok := translations["en"]; ok {
			if text, exists := trans[key]; exists {
				if len(args) > 0 {
					return fmt.Sprintf(text, args...)
				}
				return text
			}
		}
	}

	// 如果都没找到，返回key本身
	if len(args) > 0 {
		return fmt.Sprintf(key, args...)
	}
	return key
}

// 设置语言
func SetLanguage(lang string) error {
	if _, ok := supportedLanguages[lang]; !ok {
		return fmt.Errorf("不支持的语言: %s", lang)
	}

	mu.Lock()
	currentLang = lang
	if tag, ok := supportedLanguages[lang]; ok {
		printer = message.NewPrinter(tag)
	}
	mu.Unlock()

	return nil
}

// 获取当前语言
func GetCurrentLanguage() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

// 获取支持的语言列表
func GetSupportedLanguages() []string {
	var langs []string
	for lang := range supportedLanguages {
		langs = append(langs, lang)
	}
	return langs
}

// 格式化打印（支持本地化数字格式）
func Printf(format string, args ...interface{}) {
	printer.Printf(format, args...)
}

// 格式化字符串（支持本地化数字格式）
func Sprintf(format string, args ...interface{}) string {
	return printer.Sprintf(format, args...)
}
