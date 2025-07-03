package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var registryMirrors = []string{"https://docker.1ms.run", "https://docker.1panel.live"}

func linuxDockerCheck() (bool, string) {
	_, err := exec.Command("docker", "version").Output()
	if err != nil {
		rootCmd.PrintErrf("Failed to check docker version: %v", err)
		return false, "Failed to check docker version"
	}
	// 2. 检查docker是否正在运行
	dockerStatus, err := exec.Command("systemctl", "is-active", "docker").Output()
	if err != nil {
		rootCmd.PrintErrf("Failed to check docker status: %v", err)
		return false, "Failed to check docker status"
	}

	if string(dockerStatus) == "active" {
		rootCmd.Println("Docker is running")
		return true, "Running"
	} else {
		rootCmd.Println("Docker is not running")
		return true, "Stopped"
	}
}

func linuxProxy() {
	fmt.Println("🚀 开始配置 Linux Docker 镜像源...")

	// 1. 读取/etc/docker/daemon.json
	fmt.Println("📁 读取 Docker 配置文件...")
	daemonJson, err := os.ReadFile("/etc/docker/daemon.json")
	if err != nil {
		fmt.Println("❌ 读取 /etc/docker/daemon.json 失败")
		return
	}
	fmt.Println("✅ 成功读取 Docker 配置文件")

	// 2. 解析daemon.json
	fmt.Println("🔍 解析配置文件...")
	var daemonConfig map[string]interface{}
	err = json.Unmarshal(daemonJson, &daemonConfig)
	if err != nil {
		fmt.Println("❌ 解析 /etc/docker/daemon.json 失败")
		return
	}
	fmt.Println("✅ 成功解析配置文件")

	// 3. 检查是否存在registry-mirrors配置
	fmt.Println("🔍 检查现有镜像源配置...")
	mirrors, ok := daemonConfig["registry-mirrors"]
	if !ok {
		fmt.Println("📋 未发现现有镜像源配置，将创建新的配置")
		daemonConfig["registry-mirrors"] = registryMirrors
	} else {
		// 4. 检查mirrors是否为空
		if mirrors == nil {
			fmt.Println("📋 现有镜像源配置为空，将创建新的配置")
			daemonConfig["registry-mirrors"] = registryMirrors
		} else {
			// 5. 检查mirrors是否为数组
			mirrorsArray, ok := mirrors.([]interface{})
			if !ok {
				fmt.Println("⚠️  现有镜像源格式不正确，将创建新的配置")
				daemonConfig["registry-mirrors"] = registryMirrors
			} else {
				// 6. 检查mirrorsArray是否为空
				if len(mirrorsArray) == 0 {
					fmt.Println("📋 现有镜像源列表为空，将创建新的配置")
					daemonConfig["registry-mirrors"] = registryMirrors
				} else {
					fmt.Printf("📋 发现现有镜像源: %v\n", mirrorsArray)

					// 7. 追加新的镜像源
					fmt.Printf("🔄 正在添加镜像源: %v\n", registryMirrors)
					for _, mirror := range registryMirrors {
						// 检查是否已存在
						exists := false
						for _, existing := range mirrorsArray {
							if existing == mirror {
								exists = true
								break
							}
						}
						if !exists {
							mirrorsArray = append(mirrorsArray, mirror)
							fmt.Printf("✅ 添加镜像源: %s\n", mirror)
						} else {
							fmt.Printf("⏭️  镜像源已存在，跳过: %s\n", mirror)
						}
					}
					daemonConfig["registry-mirrors"] = mirrorsArray
					fmt.Printf("📝 最终镜像源列表: %v\n", mirrorsArray)
				}
			}
		}
	}

	// 8. 将daemon.json写回文件
	fmt.Println("💾 保存配置文件...")
	daemonJson, err = json.MarshalIndent(daemonConfig, "", "  ")
	if err != nil {
		fmt.Println("❌ 序列化配置文件失败")
		return
	}
	fmt.Println("✅ 成功序列化配置文件")

	// 9. 将daemon.json写回文件
	err = os.WriteFile("/etc/docker/daemon.json", daemonJson, 0644)
	if err != nil {
		fmt.Println("❌ 写入 /etc/docker/daemon.json 失败")
		return
	}
	fmt.Println("✅ 成功写入配置文件")

	// 10. 执行 systemctl daemon-reload
	fmt.Println("🔄 重新加载系统服务...")
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		fmt.Println("❌ 执行 systemctl daemon-reload 失败")
		return
	}
	fmt.Println("✅ 成功重新加载系统服务")

	// 11. 执行 systemctl restart docker
	fmt.Println("🔄 重启 Docker 服务...")
	err = exec.Command("systemctl", "restart", "docker").Run()
	if err != nil {
		fmt.Println("❌ 执行 systemctl restart docker 失败")
		return
	}
	fmt.Println("✅ 成功重启 Docker 服务")
	fmt.Println("🎉 Linux Docker 镜像源配置完成！")
}

func macDockerCheck() (bool, string) {
	// 1. 检查是否安装了orbStack
	_, err := exec.Command("orbctl", "version").Output()
	if err != nil {
		rootCmd.PrintErrln("Failed to check orbStack version")
		return false, "Failed to check orbStack version"
	}
	// 2. 检查orbStack状态
	orbStack, err := exec.Command("orbctl", "status").Output()
	if err != nil {
		return true, "Stopped"
	}
	if string(orbStack) == "Running" {
		rootCmd.Println("orbStack is running")
		return true, string(orbStack)
	} else if string(orbStack) == "Stopped" {
		rootCmd.Println("orbStack is stopped")
		return true, string(orbStack)
	}
	return true, string(orbStack)
}

func updateObrStack() {
	fmt.Println("🚀 开始配置 Docker 镜像源...")

	// 1. 获取用户主目录并构建正确的路径
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ 获取用户主目录失败:", err)
		return
	}
	dockerConfigPath := filepath.Join(homeDir, ".orbstack", "config", "docker.json")
	fmt.Printf("📁 配置文件路径: %s\n", dockerConfigPath)

	// 2. 读取~/.orbstack/config/docker.json
	dockerJson, err := os.ReadFile(dockerConfigPath)
	if err != nil {
		fmt.Printf("❌ 读取配置文件失败: %v\n", err)
		return
	}
	fmt.Println("✅ 成功读取配置文件")

	// 3. 解析现有的配置
	var config map[string]interface{}
	err = json.Unmarshal(dockerJson, &config)
	if err != nil {
		fmt.Printf("❌ 解析配置文件失败: %v\n", err)
		return
	}
	fmt.Println("✅ 成功解析配置文件")

	// 4. 获取现有的 registry-mirrors 配置
	existingMirrors, ok := config["registry-mirrors"]
	var mirrorsArray []interface{}

	if ok && existingMirrors != nil {
		// 如果存在现有的镜像配置，转换为数组
		if mirrors, ok := existingMirrors.([]interface{}); ok {
			mirrorsArray = mirrors
			fmt.Printf("📋 发现现有镜像源: %v\n", mirrorsArray)
		} else {
			fmt.Println("⚠️  现有镜像源格式不正确，将创建新的配置")
			mirrorsArray = []interface{}{}
		}
	} else {
		fmt.Println("📋 未发现现有镜像源配置，将创建新的配置")
		mirrorsArray = []interface{}{}
	}

	// 5. 追加新的镜像源
	fmt.Printf("🔄 正在添加镜像源: %v\n", registryMirrors)
	for _, mirror := range registryMirrors {
		// 检查是否已存在
		exists := false
		for _, existing := range mirrorsArray {
			if existing == mirror {
				exists = true
				break
			}
		}
		if !exists {
			mirrorsArray = append(mirrorsArray, mirror)
			fmt.Printf("✅ 添加镜像源: %s\n", mirror)
		} else {
			fmt.Printf("⏭️  镜像源已存在，跳过: %s\n", mirror)
		}
	}

	// 6. 更新配置
	config["registry-mirrors"] = mirrorsArray
	fmt.Printf("📝 最终镜像源列表: %v\n", mirrorsArray)

	// 7. 将修改后的配置转回JSON
	updatedJson, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("❌ 序列化配置失败:", err)
		return
	}
	fmt.Println("✅ 成功序列化配置")

	// 8. 写回~/.orbstack/config/docker.json
	err = os.WriteFile(dockerConfigPath, updatedJson, 0644)
	if err != nil {
		fmt.Printf("❌ 写入配置文件失败: %v\n", err)
		return
	}
	fmt.Println("✅ 成功写入配置文件")

	// 9. 重启orbStack
	fmt.Println("🔄 正在重启 Docker 服务...")
	err = exec.Command("orbctl", "restart", "docker").Run()
	if err != nil {
		fmt.Println("❌ 重启 Docker 服务失败:", err)
		return
	}
	fmt.Println("✅ 成功重启 Docker 服务")
	fmt.Println("🎉 Docker 镜像源配置完成！")
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Automatically set Docker image source address",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 判断操作系统
		os := runtime.GOOS
		if os == "linux" {
			isInstall, _ := linuxDockerCheck()
			if isInstall {
				linuxProxy()
			} else {
				rootCmd.PrintErrln("Make sure Docker is installed")
			}
		} else if os == "darwin" {
			isInstall, _ := macDockerCheck()
			if isInstall {
				updateObrStack()
			} else {
				rootCmd.PrintErrln("Make sure orbStack is installed")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
