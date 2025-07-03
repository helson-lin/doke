package main

import (
	"log"

	"github.com/helson-lin/doke/cmd"
	"github.com/helson-lin/doke/i18n"
)

func main() {
	// 初始化国际化系统
	if err := i18n.Init(); err != nil {
		log.Printf("Warning: Failed to initialize i18n: %v", err)
	}

	cmd.Execute()
}
